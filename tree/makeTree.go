package tree

import (
	c "Emissary/configs"
	"fmt"
	"strings"
)

// MakeTree - формирует дерево из массива пользователей
func MakeTree(users []c.LdapAttributes) (tree c.Branch) {

	fmt.Println("Формирую дерево пользователей по отделам")

	// Перебираем массив пользователей
	for _, user := range users {

		// Назначаем первой ветвью - ствол дерева
		currentBranch := &tree

		// Перебираем OU пользователя
		for i := 0; i < len(user.OUs); i++ {

			if _, ok := currentBranch.Branch[user.OUs[i]]; !ok {
				// Если нужной ветви нет - создаем её и назначаем текущей ("садимся на неё")
				currentBranch = addBranch(currentBranch, user.OUs[i])
			} else {
				// Иначе - назначем текущей ("садимся на неё")
				currentBranch = currentBranch.Branch[user.OUs[i]]
			}
		}
		// Когда созданы все нужные ветви - добавляем на неё должность ("листья")
		addUser(currentBranch, &user)

	}

	return tree
}

// Создает ветви на дереве
func addBranch(branch *c.Branch, ou string) *c.Branch {
	if branch.Branch == nil {
		branch.Branch = make(map[string]*c.Branch)
	}
	branch.Branch[ou] = &c.Branch{}
	return branch.Branch[ou]
}

// Добавляет пользователя на дерево
func addUser(branch *c.Branch, user *c.LdapAttributes) {
	if branch.Users == nil {
		branch.Users = make(map[string]c.BranchAttributes)
	}
	if _, ok := branch.Users[user.Name]; !ok {
		branch.Users[user.Name] = c.BranchAttributes{
			Position: user.Position,
			Contacts: strings.Join(user.Contacts, ","),
		}
	}
}
