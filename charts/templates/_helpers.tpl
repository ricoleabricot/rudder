{{/*
    Expand the name of the chart.
*/}}
{{- define "rudder.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
    Create a default fully qualified app name.
    We truncate at 63 chars because some Kubernetes
    name fields are limited to this (DNS naming spec).
    If release name contains chart name it will be
    used as a full name, otherwise, it is used as a
    suffix of the name.
*/}}
{{- define "rudder.fullname" -}}
{{- if .Values.fullnameOverride }}
    {{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
    {{- $name := default .Chart.Name .Values.nameOverride }}
    {{- if contains $name .Release.Name }}
        {{- .Release.Name | trunc 63 | trimSuffix "-" }}
    {{- else }}
        {{- printf "%s-%s" $name .Release.Name | trunc 63 | trimSuffix "-" }}
    {{- end }}
{{- end }}
{{- end }}

{{/*
    Common resource labels.
*/}}
{{- define "rudder.labels" -}}
{{- include "rudder.selectorLabels" . }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
helm.sh/chart: {{ include "rudder.chartLabel" . }}
{{- end }}

{{/*
    Pod selector labels.
*/}}
{{- define "rudder.selectorLabels" -}}
app.kubernetes.io/name: {{ include "rudder.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
   Create the name of the service account to use.
*/}}
{{- define "rudder.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
    {{- default (include "rudder.fullname" .) .Values.serviceAccount.name }}
{{- else }}
    {{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
    Chart name and version as used by the chart label.
*/}}
{{- define "rudder.chartLabel" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}
