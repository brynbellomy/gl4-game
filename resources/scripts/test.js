
function update(time, scriptInterface) {
    // console.log("time:", time)
    // var asdf = JSON.stringify(entityManager.EntitiesMatching(1 | 2 | 4))
    var entityManager = scriptInterface.EntityManager()

    var result = entityManager.MakeCmptQuery(['physics', 'position'])
    var componentQuery = result[0]
    var err = result[1]
    if (err !== null && err !== undefined) {
        throw new Error(err)
    }

    // console.log('componentQuery = ', JSON.stringify(componentQuery))
    // console.log("entityManager =", entityManager.EntitiesMatching(componentQuery))
}