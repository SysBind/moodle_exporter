# Deploy: tell Tilt what YAML to deploy
k8s_yaml('tilt/exporter.yaml')
k8s_yaml('tilt/postgres.yaml')
k8s_yaml('tilt/prometheus.yaml')
k8s_yaml('tilt/moodle.yaml')

# Build: tell Tilt what images to build from which directories
docker_build('sysbind/moodle_exporter', '.')
docker_build('sysbind/moodle-php-apache', 'tilt/', dockerfile = 'tilt/Dockerfile.moodle')


# Watch: tell Tilt how to connect locally (optional)
k8s_resource('moodle', port_forwards="8001:80", labels=["moodle"])
k8s_resource('postgresql', port_forwards="5432:5432", labels=["moodle"])
k8s_resource('prometheus', port_forwards="8000:80", labels=["prometheus"])
k8s_resource('moodle-exporter', port_forwards="2345:2345", labels=["prometheus"])
