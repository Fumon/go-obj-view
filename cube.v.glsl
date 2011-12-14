#version 120
attribute vec3 obj_coord;
attribute vec3 obj_normal;
uniform mat4 m_transform;
uniform mat3 m_3x3itModel;
varying vec3 vary_norm;

void main(void) {
	vary_norm  = normalize(m_3x3itModel * obj_normal);

	gl_Position = m_transform * vec4(obj_coord, 1.0);
}
