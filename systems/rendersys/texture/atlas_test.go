package texture_test

import (
	"os"
	"path"

	"github.com/brynbellomy/gl4-game/systems/rendersys/texture"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Atlas", func() {
	Context("when initialized from a config map", func() {
		config := map[string]interface{}{
			"name": "test-fireball",
			"animations": map[string]interface{}{
				"flying": []interface{}{
					"flying-001.png",
					"flying-002.png",
					"flying-003.png",
					"flying-004.png",
				},
			},
		}

		// assetRoot, err := getAssetPath()
		assetRoot, err := os.Getwd()
		if err != nil {
			Fail(err.Error())
		}

		assetRoot = path.Join(assetRoot, "test-fireball")

		It("should initialize all fields correctly", func() {
			atlas, err := texture.NewAtlasFromConfig(assetRoot, config)
			if err != nil {
				Fail(err.Error())
			}

			Expect(atlas.Name()).To(Equal("test-fireball"))

			anims := atlas.TextureFilenames()
			Expect(anims).To(HaveLen(1))
			Expect(anims).To(HaveKey("flying"))
			Expect(anims["flying"]).To(HaveLen(4))

			Expect(atlas.Animation("flying")).To(HaveLen(0))
		})
	})
})
