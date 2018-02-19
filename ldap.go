package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"os"

	"gopkg.in/ini.v1"
	"gopkg.in/ldap.v2"
)

type LDAPUser struct {
	DN         string
	CN         string
	Email      string
	Department string
	Division   string
}

type LdapRunner struct {
	Host           string
	Port           uint
	Tls            bool
	TlsSkipVerify  bool
	TlsServerName  string
	Username       string
	Password       string
	BaseDn         string
	connection     *ldap.Conn
	Attributes     []string
	ActiveFilter   string
	DisabledFilter string
	PageSize       uint
}

func (r *LdapRunner) connect() error {
	if r.Tls {
		err := r.connectTLS()
		if err != nil {
			return err
		}
	} else {
		err := r.connectClear()
		if err != nil {
			return err
		}
	}

	// First bind with a read only user
	err := r.connection.Bind(r.Username, r.Password)
	if err != nil {
		return err
	}
	return nil
	//	FIX defer l.Close()
}

func (r *LdapRunner) connectClear() error {
	var err error
	r.connection, err = ldap.Dial("tcp", fmt.Sprintf("%s:%d", r.Host, r.Port))
	if err != nil {
		return err
	}
	return nil
}

func (r *LdapRunner) connectTLS() error {
	var err error
	tlscfg := &tls.Config{
		InsecureSkipVerify: r.TlsSkipVerify,
		ServerName:         r.TlsServerName,
	}
	r.connection, err = ldap.DialTLS("tcp", fmt.Sprintf("%s:%d", r.Host, r.Port), tlscfg)
	if err != nil {
		return err
	}
	return nil
}

func (r *LdapRunner) close() {
	if r.connection != nil {
		r.connection.Close()
	}
}

func (r *LdapRunner) ScanForDisabledUsers() ([]LDAPUser, error) {
	searchRequest := ldap.NewSearchRequest(
		r.BaseDn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		r.DisabledFilter, r.Attributes,
		nil,
	)
	return r.Scan(searchRequest)
}

func (r *LdapRunner) ScanForActiveUsers() ([]LDAPUser, error) {
	searchRequest := ldap.NewSearchRequest(
		r.BaseDn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		r.ActiveFilter, r.Attributes,
		nil,
	)
	return r.Scan(searchRequest)
}

func (r *LdapRunner) Scan(searchRequest *ldap.SearchRequest) ([]LDAPUser, error) {

	err := r.connect()
	if err != nil {
		return nil, err
	}
	defer r.close()

	users := make([]LDAPUser, 0, 1000)
	sr, err := r.connection.SearchWithPaging(searchRequest, uint32(r.PageSize))
	if err != nil {
		return nil, err
	}
	//fmt.Printf("\nresults:\n%v\n", pretty.Formatter(sr))
	for _, entry := range sr.Entries {
		//fmt.Printf("%v: %v\n", entry.GetAttributeValue("cn"), entry.GetAttributeValue("mail"))
		users = append(users, LDAPUser{
			DN:         entry.DN,
			CN:         entry.GetAttributeValue("cn"),
			Email:      entry.GetAttributeValue("mail"),
			Department: entry.GetAttributeValue("department"),
			Division:   entry.GetAttributeValue("division"),
		})
		//entry.PrettyPrint(4)
		//fmt.Printf("%v\n", pretty.Formatter(entry))
	}

	return users, nil

}

func NewLdapRunner(cfg *ini.Section) (*LdapRunner, error) {
	password := os.Getenv("LDAP_PASSWD")
	if password == "" {
		return nil, errors.New("missing LDAP_PASSWD environment variable")
	}
	host, err := cfg.GetKey("host")
	if err != nil {
		return nil, err
	}
	port, err := cfg.Key("port").Uint()
	if err != nil {
		return nil, err
	}
	username, err := cfg.GetKey("username")
	if err != nil {
		return nil, err
	}
	basedn, err := cfg.GetKey("basedn")
	if err != nil {
		return nil, err
	}
	pagesz, err := cfg.Key("result_page_size").Uint()
	if err != nil {
		return nil, err
	}
	activeFilter, err := cfg.GetKey("activeFilter")
	if err != nil {
		return nil, err
	}
	disabledFilter, err := cfg.GetKey("disabledFilter")
	if err != nil {
		return nil, err
	}
	attr, err := cfg.GetKey("attributes")
	if err != nil {
		return nil, err
	}
	tlsFlag, err := cfg.Key("tls").Bool()
	if err != nil {
		return nil, err
	}
	tlsSkipVerify, err := cfg.Key("tls_skip_verify").Bool()
	if err != nil {
		return nil, err
	}
	tlsServerName, err := cfg.GetKey("tls_server_name")
	if err != nil {
		return nil, err
	}
	return &LdapRunner{
		Host:           host.Value(),
		Port:           port,
		Tls:            tlsFlag,
		TlsSkipVerify:  tlsSkipVerify,
		TlsServerName:  tlsServerName.Value(),
		Username:       username.Value(),
		Password:       password,
		BaseDn:         basedn.Value(),
		connection:     nil,
		Attributes:     attr.Strings(","),
		ActiveFilter:   activeFilter.Value(),
		DisabledFilter: disabledFilter.Value(),
		PageSize:       pagesz,
	}, nil
}

/*
func GetFromAD(connect *ldap.Conn, ADBaseDN, ADFilter string, ADAttribute []string, ADPage uint32) *[]LDAPElement {
	//sizelimit in searchrequest is the limit, which throws an error when the number of results exceeds the limit.
	searchRequest := ldap.NewSearchRequest(ADBaseDN, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, ADFilter, ADAttribute, nil)
	sr, err := connect.SearchWithPaging(searchRequest, ADPage)
	CheckForError(err)
	//fmt.Println(len(sr.Entries))
	ADElements := []LDAPElement{}
	for _, entry := range sr.Entries {
		NewADEntity := new(LDAPElement)
		NewADEntity.DN = entry.DN
		for _, attrib := range entry.Attributes {
			NewADEntity.attributes = append(NewADEntity.attributes, keyvalue{attrib.Name: attrib.Values})
		}
		ADElements = append(ADElements, *NewADEntity)
	}
	return &ADElements
}

func InitialrunAD(ADHost, AD_Port, ADUsername, ADPassword, ADBaseDN, ADFilter string, ADAttribute []string,
	ADPage int, ADConnTimeout int, UseTLS bool, InsecureSkipVerify bool, CRTValidFor, ADCrtPath string, shutdownChannel chan string, ADElementsChan chan *[]LDAPElement) {
	fmt.Println("Connecting to AD", ADHost)
	var connectAD *ldap.Conn
	if UseTLS == false {
		connectAD = ConnectToDirectoryServer(ADHost, AD_Port, ADUsername, ADPassword, ADConnTimeout)
	} else {
		connectAD = ConnectToDirectoryServerTLS(ADHost, AD_Port, ADUsername, ADPassword, ADConnTimeout, InsecureSkipVerify, CRTValidFor, ADCrtPath)
	}
	// defer func() {shutdownChannel <- "Done from func InitialrunAD"}()
	defer fmt.Println("closed")
	defer connectAD.Close()
	defer fmt.Println("Closing connection")
	ADElements := GetFromAD(connectAD, ADBaseDN, ADFilter, ADAttribute, uint32(ADPage))
	// fmt.Println(reflect.TypeOf(ADElements))
	fmt.Println(ADElements)
	//fmt.Println("Writing results to ", reflect.TypeOf(ADElementsChan))
	// fmt.Println("Length of ", reflect.TypeOf(ADElementsChan), "is", len(*ADElements))
	//ADElementsChan <- ADElements
	// fmt.Println("Passing", reflect.TypeOf(ADElementsChan), "to Main thread")

}
*/
