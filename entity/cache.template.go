package entity

import (
	"fmt"
	"io/ioutil"
	"sync"

	"gopkg.in/yaml.v2"

	"github.com/brynbellomy/gl4-game/systems/assetsys"
)

type (
	TemplateCache struct {
		mutex         sync.RWMutex
		templates     map[string]Entity
		fs            assetsys.IFilesystem
		entityFactory *EntityFactory
	}
)

func NewTemplateCache(fs assetsys.IFilesystem, entityFactory *EntityFactory) *TemplateCache {
	return &TemplateCache{
		mutex:         sync.RWMutex{},
		templates:     map[string]Entity{},
		fs:            fs,
		entityFactory: entityFactory,
	}
}

func (c *TemplateCache) Load(tplname string) (Entity, error) {
	tpl, err := c.loadFromCache(tplname)
	if err != nil {
		return Entity{}, err
	}

	// clone all of the template's components
	cmpts := make([]IComponent, len(tpl.Components))
	for i, cmpt := range tpl.Components {
		cmpts[i] = cmpt.Clone()
	}

	tpl.Components = cmpts

	return tpl, nil
}

func (c *TemplateCache) loadFromCache(tplname string) (Entity, error) {
	c.mutex.RLock()
	template, exists := c.templates[tplname]
	c.mutex.RUnlock()

	if exists {
		return template, nil
	}

	fmt.Println("entity cache: loading", tplname)

	c.mutex.Lock()
	defer c.mutex.Unlock()

	template, err := c.loadTemplateFile(tplname)
	if err != nil {
		return Entity{}, err
	}

	c.templates[tplname] = template
	return template, nil
}

func (c *TemplateCache) loadTemplateFile(tplname string) (Entity, error) {
	f, err := c.fs.OpenFile(tplname+".yaml", 0, 0400)
	if err != nil {
		return Entity{}, err
	}

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return Entity{}, err
	}

	var cfg map[string]interface{}
	err = yaml.Unmarshal(bytes, &cfg)
	if err != nil {
		return Entity{}, err
	}

	entity, err := c.entityFactory.EntityFromConfig(cfg)
	return entity, err
}
