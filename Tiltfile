# Deploy: tell Tilt what YAML to deploy
k8s_yaml('tilt/exporter.yaml')
k8s_yaml('tilt/postgres.yaml')
# Build: tell Tilt what images to build from which directories

docker_build('sysbind/moodle_exporter', '.')


# Watch: tell Tilt how to connect locally (optional)

# k8s_resource('api', port_forwards="5734:5000", labels=["backend"])
