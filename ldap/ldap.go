package ldap

import (
	"fmt"
	"strings"

	c "Emissary/configs"

	ldap "github.com/go-ldap/ldap/v3"
)

// Список возвращаемых атрибутов
var requestAttributes = []string{
	"description",
	"telephoneNumber",
	"otherTelephone",
	"mobile",
	"otherMobile",
	"mail",
}

// GetUsersFromLdap - запрашивает пользователей из LDAP и возвращает их список
func GetUsersFromLdap() (users []c.LdapAttributes, err error) {
	fmt.Println("Подключаюсь к LDAP")

	// Формирование подключения к ldap
	l, err := ldap.Dial("tcp", c.LdapServer)
	if err != nil {
		return nil, err
	}
	defer l.Close()

	// Использование УЗ для подключения
	err = l.Bind(
		// Проверка наличия домена LDAP. Если он есть, использовать с логином
		func() string {
			if c.LdapDomain == "" {
				return c.Login
			} else {
				return c.Login + "@" + c.LdapDomain
			}
		}(),

		c.Pass)
	if err != nil {
		return nil, err
	}

	// Создание запроса в AD
	searchRequest := ldap.NewSearchRequest(
		c.LdapBaseDn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		c.LdapFilter,
		requestAttributes,
		nil,
	)

	fmt.Println("Запрашиваю пользователей")

	// Выполнение запроса в AD
	sr, err := l.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Вернулось %v пользовател(я/ей)\n", len(sr.Entries))

	// Инициализирую массив пользователей и складываю в него данные
	users = make([]c.LdapAttributes, 0)

	// Прохожусь по всем сущностям и формирую список пользователей
	for _, entry := range sr.Entries {
		var user c.LdapAttributes

		// Должность пользователя
		user.Position = entry.GetAttributeValue("description")

		// Суммирую все контактные данные пользователя
		for _, attr := range requestAttributes {
			user.Contacts = append(user.Contacts, entry.GetAttributeValues(attr)...)
		}

		// Разбираю DN строку для получения имени и каталогов пользователя
		dn := strings.TrimSuffix(entry.DN, ","+c.LdapBaseDn) // Отрезаю BaseDn
		for strings.Contains(dn, "OU=") {

			indexOU := strings.LastIndex(dn, ",OU=")    // Расположение последнего каталога в строке
			user.OUs = append(user.OUs, dn[indexOU+4:]) // Прибавляю каталог к массиву каталогов пользователя
			dn = dn[:indexOU]                           // Отрезаю данный каталог от строки
		}

		user.Name = strings.TrimPrefix(dn, "CN=") // Получаю имя пользователя

		users = append(users, user)
	}
	return users, nil
}
