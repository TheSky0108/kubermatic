#!/usr/bin/env bash

# Copyright 2020 The Kubermatic Kubernetes Platform contributors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -euo pipefail

cd $(dirname $0)/../..
source hack/lib.sh

apk add --no-cache -U git bash openssh
git fetch --tags

# we only synchronize charts for the master branch and tagged versions
if [ "$PULL_BASE_REF" != "master" ] && ! (git tag -l | grep -q "^$PULL_BASE_REF\$"); then
  echo "Base ref $PULL_BASE_REF is neither master nor a tag! Exiting..."
  exit 1
fi

# The dashboard tag should match the Kubermatic tag for tagged releases
KUBERMATIC_VERSION="$PULL_BASE_REF"
DASHBOARD_VERSION="$KUBERMATIC_VERSION"

# find out which kubermatic-installer branch to synchronize to
if [ "$PULL_BASE_REF" == "master" ]; then
  INSTALLER_BRANCH="master"
  # do not use the actual head hash to avoid spamming the installer
  # repo for every new commit on master branch
  KUBERMATIC_VERSION="latest"
  DASHBOARD_VERSION="master"
elif grep -qE '^v[0-9]+\.[0-9]+\.[0-9]+.*' <<< "$PULL_BASE_REF"; then
  # e.g. 'release/v2.12'
  INSTALLER_BRANCH="release/$(grep -o -E '^v[0-9]+\.[0-9]+' <<< "$PULL_BASE_REF")"
elif grep -E '^weekly-[0-9]{4}-.*' <<< "$PULL_BASE_REF"; then
  INSTALLER_BRANCH="weekly"
else
  echo "I don't know to which installer branch the tag '$PULL_BASE_REF' belongs to. Exiting..."
  exit 1
fi

sed -i "s/__KUBERMATIC_TAG__/$KUBERMATIC_VERSION/g" charts/*/*.yaml
sed -i "s/__DASHBOARD_TAG__/$DASHBOARD_VERSION/g" charts/*/*.yaml

git config --global user.email "dev@loodse.com"
git config --global user.name "Prow CI Robot"
git config --global core.sshCommand 'ssh -o CheckHostIP=no -i /ssh/id_rsa'
ensure_github_host_pubkey

export CHARTS='kubermatic kubermatic-operator cert-manager certs nginx-ingress-controller nodeport-proxy oauth minio iap s3-exporter'
export MONITORING_CHARTS='alertmanager blackbox-exporter grafana kube-state-metrics node-exporter prometheus'
export LOGGING_CHARTS='loki promtail elasticsearch kibana fluentbit'
export BACKUP_CHARTS='velero'
export CHARTS_DIR=$(pwd)/charts
export TARGET_DIR='sync_target'
export TARGET_VALUES_FILE=${TARGET_DIR}/values.example.yaml
export TARGET_VALUES_SEED_FILE=${TARGET_DIR}/values.seed.example.yaml

# create fresh clone of the installer repository
rm -rf ${TARGET_DIR}
mkdir ${TARGET_DIR}
git clone git@github.com:kubermatic/kubermatic-installer.git ${TARGET_DIR}
cd ${TARGET_DIR}
git checkout ${INSTALLER_BRANCH} || git checkout -b ${INSTALLER_BRANCH}
cd ..

# re-assemble example values.yaml
rm -f ${TARGET_DIR}/values.yaml

for VALUE_FILE in ${TARGET_VALUES_FILE} ${TARGET_VALUES_SEED_FILE}; do
  rm -f ${VALUE_FILE}
  echo "# THIS FILE IS GENERATED BY https://github.com/kubermatic/kubermatic/blob/master/hack/ci/ci-sync-charts.sh" > ${VALUE_FILE}
done

cat "${CHARTS_DIR}/kubermatic/values.yaml" >> ${TARGET_VALUES_SEED_FILE}

# ensure that charts we don't know yet (from a future version)
# also get cleaned up
rm -rf ${TARGET_DIR}/charts
mkdir ${TARGET_DIR}/charts

# sync base charts
for CHART in ${CHARTS}; do
  echo "syncing ${CHART}..."
  cp -r ${CHARTS_DIR}/${CHART} ${TARGET_DIR}/charts/${CHART}

  echo "# ====== ${CHART} ======" >> ${TARGET_VALUES_FILE}
  cat "${CHARTS_DIR}/${CHART}/values.yaml" >> ${TARGET_VALUES_FILE}
  echo "" >> ${TARGET_VALUES_FILE}
done

# sync monitoring charts
echo "" >> ${TARGET_VALUES_FILE}
echo "# ========================" >> ${TARGET_VALUES_FILE}
echo "# ====== Monitoring ======" >> ${TARGET_VALUES_FILE}
echo "# ========================" >> ${TARGET_VALUES_FILE}
echo "" >> ${TARGET_VALUES_FILE}
mkdir -p "${TARGET_DIR}/charts/monitoring"
for CHART in ${MONITORING_CHARTS}; do
  echo "syncing ${CHART}..."
  cp -r ${CHARTS_DIR}/monitoring/${CHART} ${TARGET_DIR}/charts/monitoring/${CHART}

  echo "# ====== ${CHART} ======" >> ${TARGET_VALUES_FILE}
  cat "${CHARTS_DIR}/monitoring/${CHART}/values.yaml" >> ${TARGET_VALUES_FILE}
  echo "" >> ${TARGET_VALUES_FILE}
done

# sync logging charts
echo "" >> ${TARGET_VALUES_FILE}
echo "# =======================" >> ${TARGET_VALUES_FILE}
echo "# ======= Logging =======" >> ${TARGET_VALUES_FILE}
echo "# =======================" >> ${TARGET_VALUES_FILE}
echo "" >> ${TARGET_VALUES_FILE}
mkdir -p "${TARGET_DIR}/charts/logging"
for CHART in ${LOGGING_CHARTS}; do
  echo "syncing ${CHART}..."
  cp -r ${CHARTS_DIR}/logging/${CHART} ${TARGET_DIR}/charts/logging/${CHART}

  echo "# ====== ${CHART} ======" >> ${TARGET_VALUES_FILE}
  cat "${CHARTS_DIR}/logging/${CHART}/values.yaml" >> ${TARGET_VALUES_FILE}
  echo "" >> ${TARGET_VALUES_FILE}
done

# sync backup charts
echo "" >> ${TARGET_VALUES_FILE}
echo "# =======================" >> ${TARGET_VALUES_FILE}
echo "# ======= Backups =======" >> ${TARGET_VALUES_FILE}
echo "# =======================" >> ${TARGET_VALUES_FILE}
echo "" >> ${TARGET_VALUES_FILE}
mkdir -p "${TARGET_DIR}/charts/backup"
for CHART in ${BACKUP_CHARTS}; do
  echo "syncing ${CHART}..."
  cp -r ${CHARTS_DIR}/backup/${CHART} ${TARGET_DIR}/charts/backup/${CHART}

  echo "# ====== ${CHART} ======" >> ${TARGET_VALUES_FILE}
  cat "${CHARTS_DIR}/backup/${CHART}/values.yaml" >> ${TARGET_VALUES_FILE}
  echo "" >> ${TARGET_VALUES_FILE}
done

# merge duplicate top-level keys
KNOWN_KEYS=()

while IFS= read LINE; do
  MATCH=$(echo "${LINE}" | grep -oE "^[a-zA-Z0-9]+" || true)

  if [ -z "${MATCH}" ]; then
    echo "${LINE}"
  else
    if ! [[ " ${KNOWN_KEYS[@]} " =~ " ${MATCH} " ]]; then
      KNOWN_KEYS+=("${MATCH}")
      echo "${LINE}"
    fi
  fi
done < ${TARGET_VALUES_FILE} > ${TARGET_DIR}/values.example.tmp.yaml

mv ${TARGET_DIR}/values.example.{tmp.,}yaml

# assemble static manifests to make installing the Kubermatic Operator easier
mkdir -p ${TARGET_DIR}/manifests

CRDS_MANIFEST=${TARGET_DIR}/manifests/kubermatic-crds.yaml
cat << EOF > ${CRDS_MANIFEST}
# Kubermatic $PULL_BASE_REF CRDs

EOF

for file in ${CHARTS_DIR}/kubermatic/crd/*.yaml; do
  cat "$file" >> ${CRDS_MANIFEST}
  echo -e "\n---" >> ${CRDS_MANIFEST}
done

# commit and push
cd ${TARGET_DIR}
git add .
if ! git status|grep 'nothing to commit'; then
  if [ "$INSTALLER_BRANCH" == "master" ]; then
    SHORT_HASH="$(git rev-parse --short HEAD | tr -d '\n')"
    git commit -m "Syncing charts from master branch @ $SHORT_HASH"
  else
    git commit -m "Syncing charts from release ${PULL_BASE_REF}"
    # $PULL_BASE_REF is a tag, we've checked that earlier
    git tag $PULL_BASE_REF
  fi

  git push --tags origin ${INSTALLER_BRANCH}
fi

cd ..
rm -rf ${TARGET_DIR}
