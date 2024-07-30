package temps

import (
	"os"
	"text/template"
)

func RtmpNginxFrame() {
	// ####################################################
	//  rabbit template
	rtmp_tmpl, err := template.New("RenderData").Parse(ngnixDockerTemplate)
	if err != nil {
		panic(err)
	}

	// Create the models directory if it does not exist
	err = os.MkdirAll("realtime", os.ModePerm)
	if err != nil {
		panic(err)
	}

	rtmp_file, err := os.Create("realtime/rtmp.Dockerfile")
	if err != nil {
		panic(err)
	}
	defer rtmp_file.Close()

	err = rtmp_tmpl.Execute(rtmp_file, RenderData)
	if err != nil {
		panic(err)
	}
}

func NginxConfFrame() {
	// ####################################################
	//  rabbit template
	rtmp_tmpl, err := template.New("RenderData").Parse(nginxConfTemplate)
	if err != nil {
		panic(err)
	}

	// Create the models directory if it does not exist
	err = os.MkdirAll("realtime", os.ModePerm)
	if err != nil {
		panic(err)
	}

	rtmp_file, err := os.Create("realtime/nginx.conf")
	if err != nil {
		panic(err)
	}
	defer rtmp_file.Close()

	err = rtmp_tmpl.Execute(rtmp_file, RenderData)
	if err != nil {
		panic(err)
	}
}

var nginxConfTemplate = `
worker_processes auto;
rtmp_auto_push on;
events {}
rtmp {
    server {
        listen 1935;
        listen [::]:1935 ipv6only=on;
 
        application live {
            live on;
            record off;
            hls on;
            hls_path /tmp/hls;
        }
    }
}

http {

    server {

        listen      8080;   

        location /hls {
            # Serve HLS fragments
            types {
                application/vnd.apple.mpegurl m3u8;
                video/mp2t ts;
            }
            root /tmp;
            add_header Cache-Control no-cache;
            add_header Access-Control-Allow-Origin *; 
        }

        location /dash {
            # Serve DASH fragments
            root /tmp;
            add_header Cache-Control no-cache;
        }
    }
}

`

var ngnixDockerTemplate = `
FROM tiangolo/nginx-rtmp
 
COPY nginx.conf /etc/nginx/nginx.conf
`
