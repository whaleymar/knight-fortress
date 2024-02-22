#version 330

uniform sampler2D texture1;
// uniform sampler2D texture2;
uniform float millis;

in vec2 fragTexCoord;

out vec4 outputColor;

float invert(float x) {
    return 1-x;
}

float norm(float x) {
    return (x + 1.)/2.;
}

vec2 norm(vec2 vec) {
    return vec2(norm(vec.x), norm(vec.y));
}

void main() {
    outputColor = texture(texture1, fragTexCoord);

    // this is a way to render textures with binary alpha values without sorting them, but it uses branching so probably slow
    // vec4 texColor = texture(texture1, fragTexCoord);
    // if (texColor.a < 0.01)
    //     discard; // works
    // outputColor = vec4(texColor.x, texColor.y, texColor.z, clampAlpha(texColor.a)); // doesn't work
}
