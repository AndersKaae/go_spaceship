// circle_mask.fs
#version 330

in vec2 fragTexCoord;
in vec4 fragColor;

uniform sampler2D texture0;
uniform vec2 center;
uniform float radius;

out vec4 finalColor;

void main() {
    vec2 uv = fragTexCoord;
    vec2 pixelPos = uv * vec2(textureSize(texture0, 0));
    float dist = distance(pixelPos, center);
    
    if (dist > radius) {
        discard;  // Ignore pixels outside circle
    }

    finalColor = texture(texture0, uv) * fragColor;
}

