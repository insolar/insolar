// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package metrics

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"time"

	"github.com/insolar/insolar/version"
)

var statusTmpl = `
<html>
<head>
<title>_status page</title>
<style>
	* {
		box-sizing: border-box;
	}

	h1 {
		text-transform: uppercase;
	}

	section {
		width: 100%;
		border: 1px solid blue;
	}

	section header {
		text-transform:uppercase;
		width: 100%;
		height: 48px;
		padding: 0 12px;
		line-height: 48px;
		font-size: 24px;
		background-color:blue;
		color: whitesmoke;
	}

	dl {
		display: grid;
		grid-template-columns: 200px 1fr;
		grid-template-rows: 1fr;
	}

	dl dt {
		grid-column: 1;
		font-weight: bold;
		text-transform: capitalize;
		text-align: right;
	}

	dl dd {
		grid-column: 2;
	}
</style>
</head>
<body>

<h1>STATUS</h1>

<h2>Build info</h2>
<pre>
{{ .VersionInfo }}
</pre>

<section>
<header>General</header>
<div class="content">
<dl>
<dt>Uptime:</dt> <dd>{{ .Uptime }}</dd>
<dt>metrics:</dt> <dd><a href="/metrics">/metrics</a></dd>
<dt>pprof:</dt> <dd><a href="/debug/pprof">/debug/pprof</a></dd>
<dt>rpcz:</dt> <dd> <a href="/debug/rpcz">/debug/rpcz</a></dd>
<dt>tracez:</dt> <dd><a href="/debug/tracez">/debug/tracez</a></dd>
</dl>
</div>
</section>

</body>
</html>
`

var parsedStatusTmpl = template.Must(template.New("proc_status").Parse(statusTmpl))

type procStatus struct {
	StartTime time.Time
}

func newProcStatus() *procStatus {
	info := &procStatus{
		StartTime: time.Now(),
	}
	return info
}

func (ps *procStatus) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var b bytes.Buffer
	err := parsedStatusTmpl.Execute(&b, struct {
		VersionInfo string
		Uptime      string
	}{
		VersionInfo: version.GetFullVersion(),
		Uptime:      fmt.Sprintf("%v", time.Since(ps.StartTime)),
	})
	if err != nil {
		http.Error(w, fmt.Sprintln("Template error:", err),
			http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	_, err = io.Copy(w, &b)
	if err != nil {
		http.Error(w, fmt.Sprintln("Copy error:", err),
			http.StatusInternalServerError)
		return
	}
}
