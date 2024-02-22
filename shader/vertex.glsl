#version 330

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

uniform vec3 offset;

in vec3 vert;
in vec2 vertTexCoord;

out vec2 fragTexCoord;

void main() {
    // only use x and y from offset because 2D
    vec4 totalOffset = vec4(offset.x, offset.y, 0., 0.) + vec4(1., 1., 0., 0.); // my opengl has -1,-1 as the bottom left of the screen
    gl_Position = projection * camera * vec4(vert, 1.) + totalOffset;
    
    fragTexCoord = vertTexCoord;
}
