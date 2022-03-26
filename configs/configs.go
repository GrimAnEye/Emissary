package configs

type (
	// ldapAttributes - атрибуты пользователя в AD
	LdapAttributes struct {
		Name     string   // Имя пользователя
		Position string   // Должность пользователя
		OUs      []string // Каталоги расположения пользователя
		Contacts []string // Контактные данные
	}

	// branchAttributes - атрибуты пользователя для вывода
	BranchAttributes struct {
		Position string // Должность пользователя
		Contacts string // Контактные данные пользователя
	}

	// branch - структура для составления дерева пользователей, в зависимости от их расположения по отделам
	Branch struct {
		Users  map[string]BranchAttributes // Пользователи на ветке
		Branch map[string]*Branch          // Подветви текущей ветки
	}
)

var (
	// Запрашиваемые переменные от пользователя

	Login      string // Имя пользователя для доступа к LDAP
	Pass       string // Пароль пользователя для доступа к LDAP
	LdapDomain string // Домен для логина пользователя (ex.example.com)
	LdapServer string // Адрес LDAP сервера
	LdapBaseDn string // Базовый каталог LDAP для поиска
	LdapFilter string // LDAP фильтр для поиска пользователей
	Output     string // Путь до файла вывода телефонного справочника

)

const (
	// Описание запрашиваемых переменных

	CLogin      string = "Имя пользователя для доступа к LDAP"
	CPass       string = "Пароль пользователя для доступа к LDAP"
	CLdapDomain string = "Домен для логина пользователя: ex.example.com"
	CLdapServer string = "Адрес LDAP сервера в формате: ldap.example.com:389"
	CLdapBaseDn string = "Базовый каталог LDAP для поиска: OU=Users,DC=exa,DC=example,DC=com"
	CLdapFilter string = "LDAP фильтр для поиска пользователей: (objectClass=organizationalPerson)"
	COutput     string = "Путь до файла вывода телефонного справочника"

	CUsage string = `Emissary - программа офисной техподдержки,
	формирующая HTML страницу "телефонного" справочника пользователей
	с их контактными данными, для offline использования`
	CErr string = `Требуется указать данные для всех ключей программы! Справка "-h"`
)
