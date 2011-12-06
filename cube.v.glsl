#version 120
attribute vec3 coord3d;
uniform mat4 m_transform;
void main(void) {
	gl_Position = m_transform * vec4(coord3d, 1.0);
}
