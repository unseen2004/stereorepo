#version 330 core

layout(location=0) out vec4 color;

uniform vec3 textColor;

void main(){
    color = vec4(textColor, 1.0);
}
