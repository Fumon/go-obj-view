#version 120
attribute vec3 obj_coord;
attribute vec3 obj_normal;
uniform mat4 m_transform;
uniform mat3 m_3x3itModel;
varying vec3 vary_norm;
varying vec4 vary_pos;

void main(void) {
	vary_norm  = normalize(m_3x3itModel * obj_normal);
	vary_pos = m_transform * vec4(obj_coord, 1.0);

	gl_Position = m_transform * vec4(obj_coord, 1.0);
}
