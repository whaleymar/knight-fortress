#version 330

uniform sampler2D tex;
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
    outputColor = texture(tex, fragTexCoord);
    // outputColor = vec4(fract(fragTexCoord), 1., 1.);
    // outputColor = texture(tex, norm(fragTexCoord));
    // outputColor = texture(tex, vec2((fragTexCoord.x + 2.)/4.), (fragTexCoord.y + 2.)/4.);
    // outputColor = vec4(1.);
}
