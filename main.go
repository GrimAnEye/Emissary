package main

import (
	c "Emissary/configs"
	l "Emissary/ldap"
	t "Emissary/tree"
	"flag"
	"fmt"
	"os"
)

func init() {

	// Инициализирую ключи программы
	flag.StringVar(&c.Login, "l", "", c.CLogin)
	flag.StringVar(&c.Pass, "p", "", c.CPass)
	flag.StringVar(&c.LdapDomain, "d", "", c.CLdapDomain)
	flag.StringVar(&c.LdapServer, "s", "", c.CLdapServer)
	flag.StringVar(&c.LdapBaseDn, "b", "", c.CLdapBaseDn)
	flag.StringVar(&c.LdapFilter, "f", "", c.CLdapFilter)
	flag.StringVar(&c.Output, "o", "Emissary.html", c.COutput)

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), c.CUsage+"\n\n")
		flag.PrintDefaults()
	}
}

func main() {

	// Запрашиваю аргументы командной строки и проверяю их наличие
	flag.Parse()
	if c.Login == "" || c.Pass == "" ||
		c.LdapServer == "" || c.LdapBaseDn == "" || c.LdapFilter == "" {
		flag.Usage()
		panic(c.CErr)
	}

	// Запрашиваю пользователей из LDAP
	users, err := l.GetUsersFromLdap()
	if err != nil {
		fmt.Printf("При запросе пользователей в LDAP произошла ошибка:\n%s\n", err)
		os.Exit(1)
	}

	// Формирую дерево пользователей по отделам
	tree := t.MakeTree(users)

	// Вывожу дерево пользователей в файл
	if err := t.PrintTree(&tree); err != nil {
		fmt.Printf("При создании файла вывода справочника возникла ошибка:\n%s\n", err)
		os.Exit(1)
	}

}
