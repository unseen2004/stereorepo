#version 330 core

layout(location=0) out vec4 color;

in DATA{
    vec2 tc;
    vec3 position;
} fs_in;

uniform sampler2D tex;
uniform int top;
uniform vec2 bird;

void main(){
    vec2 coord = fs_in.tc;
    if(top == 1){
        coord.y = 1.0 - coord.y;
    }
    color = texture(tex, coord);
    if(color.w < 1.0){
        discard;
    }
    color *= 3.0 / (length(bird - fs_in.position.xy) + 2.5) + 0.3;
    color.w = 1.0;

}