#version 330 core

layout(location=0) out vec4 color;

uniform float alpha;

void main(){
    color = vec4(0.0, 0.0, 0.0, alpha);
}

