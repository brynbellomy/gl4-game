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
		templates     map[string][]IComponent
		fs            assetsys.IFilesystem
		entityFactory *EntityFactory
	}
)

func NewTemplateCache(fs assetsys.IFilesystem, entityFactory *EntityFactory) *TemplateCache {
	return &TemplateCache{
		mutex:         sync.RWMutex{},
		templates:     map[string][]IComponent{},
		fs:            fs,
		entityFactory: entityFactory,
	}
}

func (c *TemplateCache) Load(tplname string) ([]IComponent, error) {
	tpl, err := c.loadFromCache(tplname)
	if err != nil {
		return nil, err
	}

	// clone all of the template's components
	cmpts := make([]IComponent, len(tpl))
	for i, cmpt := range tpl {
		cmpts[i] = cmpt.Clone()
	}

	return cmpts, nil
}

func (c *TemplateCache) loadFromCache(tplname string) ([]IComponent, error) {
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
		return nil, err
	}

	c.templates[tplname] = template
	return template, nil
}

func (c *TemplateCache) loadTemplateFile(tplname string) ([]IComponent, error) {
	f, err := c.fs.OpenFile(tplname+".yaml", 0, 0400)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var cfg map[string]interface{}
	err = yaml.Unmarshal(bytes, &cfg)
	if err != nil {
		return nil, err
	}

	_, cmpts, err := c.entityFactory.EntityFromConfig(cfg)
	return cmpts, err
}
