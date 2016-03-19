package scriptsys

// import (
//     "github.com/yuin/gopher-lua"

//     "github.com/brynbellomy/gl4-game/entity"
// )

// const luaEntityManagerTypeName = "entityManager"

// // Registers my entityManager type to given L.
// func registerEntityManagerType(L *lua.LState) {
//     mt := L.NewTypeMetatable(luaEntityManagerTypeName)
//     L.SetGlobal("entityManager", mt)

//     // static attributes
//     // L.SetField(mt, "new", L.NewFunction(newEntityManager))

//     // methods
//     L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
//         // "name": entityManagerGetSetName,
//         "makeCmptQuery": entityManagerMakeCmptQuery,
//     }))
// }

// func entityManagerMakeCmptQuery(L *lua.LState) int {
//     em := checkEntityManager(L)
//     if L.GetTop() != 2 {
//         L.ArgError(2, "wrong number of arguments")
//         return 0
//     }

//     mask := L.CheckInt64(2)

//     query, err := em.MakeCmptQuery(entity.ComponentMask(uint64(mask)))
//     if err != nil {
//         L.RaiseError(err)
//         return 0
//     }

//     L.Push(lua.LNumber(query))
//     return 1
// }

// // Getter and setter for the EntityManager.Name
// // func entityManagerGetSetName(L *lua.LState) int {
// //     p := checkEntityManager(L)
// //     if L.GetTop() == 2 {
// //         p.Name = L.CheckString(2)
// //         return 0
// //     }
// //     L.Push(lua.LString(p.Name))
// //     return 1
// // }

// // Checks whether the first lua argument is a *LUserData with *EntityManager and returns this *EntityManager.
// func checkEntityManager(L *lua.LState) *entity.Manager {
//     ud := L.CheckUserData(1)
//     if v, ok := ud.Value.(*entity.Manager); ok {
//         return v
//     }
//     L.ArgError(1, "entityManager expected")
//     return nil
// }