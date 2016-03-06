package mainscene

import (
	"path"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/common"
	"github.com/brynbellomy/gl4-game/entity"
	"github.com/brynbellomy/gl4-game/systems/animationsys"
	"github.com/brynbellomy/gl4-game/systems/movesys"
	"github.com/brynbellomy/gl4-game/systems/physicssys"
	"github.com/brynbellomy/gl4-game/systems/positionsys"
	"github.com/brynbellomy/gl4-game/systems/projectilesys"
	"github.com/brynbellomy/gl4-game/systems/rendersys"
	"github.com/brynbellomy/gl4-game/systems/rendersys/texture"
)

type FireballFactory struct {
	AssetRoot string

	fireballTexture uint32
	fireballAtlas   *texture.Atlas
}

func NewFireballFactory(assetRoot string) (*FireballFactory, error) {
	f := &FireballFactory{
		AssetRoot: assetRoot,
	}

	err := f.init()
	return f, err
}

func (f *FireballFactory) init() error {
	fireballTexture, err := texture.Load(path.Join(f.AssetRoot, "textures/fireball/flying-001.png"))
	if err != nil {
		return err
	}

	f.fireballTexture = fireballTexture

	f.fireballAtlas = texture.NewAtlas()
	err = f.fireballAtlas.LoadAnimation("flying", []string{
		path.Join(f.AssetRoot, "textures/fireball/flying-001.png"),
		path.Join(f.AssetRoot, "textures/fireball/flying-002.png"),
		path.Join(f.AssetRoot, "textures/fireball/flying-003.png"),
		path.Join(f.AssetRoot, "textures/fireball/flying-004.png"),
	})
	return err
}

func (f *FireballFactory) Build(pos mgl32.Vec2, vec mgl32.Vec2) ([]entity.IComponent, error) {
	boundingBox := physicssys.BoundingBox{
		{-0.1, -0.07},
		{-0.1, 0.07},
		{0.1, 0.07},
		{0.1, -0.07},
	}

	var (
		vertexShader = `
            #version 410

            uniform mat4 projection;
            uniform mat4 camera;
            uniform mat4 model;

            in vec3 vert;
            in vec2 vertTexCoord;

            out vec2 fragTexCoord;

            void main() {
                fragTexCoord = vertTexCoord;
                gl_Position = projection * camera * model * vec4(vert, 1);
            }
        ` + "\x00"

		fragmentShader = `
            #version 410

            uniform sampler2D tex;

            in vec2 fragTexCoord;

            out vec4 outputColor;

            #define M_PI 3.1415926535897932384626433832795

            vec3 rgb2hsv(vec3 c) {
                vec4 K = vec4(0.0, -1.0 / 3.0, 2.0 / 3.0, -1.0);
                vec4 p = mix(vec4(c.bg, K.wz), vec4(c.gb, K.xy), step(c.b, c.g));
                vec4 q = mix(vec4(p.xyw, c.r), vec4(c.r, p.yzx), step(p.x, c.r));

                float d = q.x - min(q.w, q.y);
                float e = 1.0e-10;
                return vec3(abs(q.z + (q.w - q.y) / (6.0 * d + e)), d / (q.x + e), q.x);
            }

            vec3 hsv2rgb(vec3 c) {
                vec4 K = vec4(1.0, 2.0 / 3.0, 1.0 / 3.0, 3.0);
                vec3 p = abs(fract(c.xxx + K.xyz) * 6.0 - K.www);
                return c.z * mix(K.xxx, clamp(p - K.xxx, 0.0, 1.0), c.y);
            }

            void main() {
                vec3 glowColor = vec3(1, 0.8, 0.8);
                vec4 texColor  = texture(tex, fragTexCoord);

                float sat = sin(M_PI * fragTexCoord.x) * sin(M_PI * fragTexCoord.y);

                vec3 mixed = mix(texColor.rgb, glowColor, 0.3 * sat);

                // vec3 texHSV = rgb2hsv(mixed);
                // texHSV.y = clamp(0.3 + sat * texHSV.y, 0.3, 1.0);
                // texHSV.z = clamp(0.3 + sat * texHSV.z, 0.6, 1.0);
                outputColor = vec4(
                    // hsv2rgb(texHSV),
                    mixed,
                    texColor.a
                );
            }

        ` + "\x00"

		spriteNode = rendersys.NewSpriteNode(vertexShader, fragmentShader)
	)

	return []entity.IComponent{
		positionsys.NewComponent(pos, common.Size{0.2, 0.14}, 2, 0),
		physicssys.NewComponent(mgl32.Vec2{}, 8, mgl32.Vec2{}, boundingBox, uint64(ProjectileCollider), uint64(EnemyCollider)),
		rendersys.NewComponent(spriteNode, f.fireballTexture),
		animationsys.NewComponent(f.fireballAtlas, "flying", true, 0, 12),
		movesys.NewComponent(mgl32.Vec2{0, 0}),
		projectilesys.NewComponent(vec, 0.01, 10, projectilesys.Firing, true),
	}, nil
}
