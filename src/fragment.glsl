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

void main() {
    outputColor = texture(tex, fragTexCoord);
    // outputColor = vec4(1.);

    // fragTextCoord is in range (-1,1)
    // outputColor = vec4((fragTexCoord.x + 1)/2., (fragTexCoord.y+1)/2., 1., 1.);

    // fract
    // vec2 newPos = fract(fragTexCoord * 10.);
    // outputColor = vec4((newPos.x + 1)/2., (newPos.y+1)/2., 1., 1.);

    // mix
 //    vec4 c1 = vec4(0.5, 0.1, 0.9, 1.);
	// vec4 c2 = vec4(0.1, 0.8, 0.7, 1.);
	// vec4 c = mix(c1, c2, fragTexCoord.x);
 //    outputColor = c;

    // sin -- in -1, 1 range, need to clamp
    // float c = sin(fragTexCoord.x * 16.); // bad
    // float c = sin((fragTexCoord.x+1)/2. * 16.); // better
    // float c = (sin((fragTexCoord.x+1)/2. * 16.) + 1)/2.; // smoothest
    // float c = (sin((fragTexCoord.x+1)/2. * 16. + millis) + 1)/2.; // animated
    // outputColor = vec4(c, 0., 1., 1.);
    
    // SDFs to draw circles
    // vec3 circle = vec3(0., 0., 0.5); // (xcenter, ycenter, radius)
    // float d = length(fragTexCoord - circle.xy) - circle.z;
    // // d = step(0., d);
    // d = invert(smoothstep(0., 0.01, d));
    // vec2 newpos = vec2(norm(fragTexCoord.x), norm(fragTexCoord.y));
    //
    // outputColor = vec4(d * newpos.x, d * newpos.y, d, 1.);
}
