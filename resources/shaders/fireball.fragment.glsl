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