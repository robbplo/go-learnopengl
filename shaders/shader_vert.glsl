#version 330 core
layout (location = 0) in vec3 vertexPosition_modelspace;
layout (location = 1) in vec3 vertexColor;
out vec3 fragmentColor;

uniform mat4 model; // Model matrix
uniform mat4 view; // View matrix
uniform mat4 projection; // Projection matrix
uniform float angle; // Rotation angle


void main()
{
    fragmentColor = vertexColor;
  
    mat4 rotation = mat4(
        cos(angle), -sin(angle), 0, 0,
        sin(angle), cos(angle), 0, 0,
        0, 0, 1, 0,
        0, 0, 0, 1
    );
    gl_Position = projection * view * model * rotation * vec4(vertexPosition_modelspace, 1.0);
}
