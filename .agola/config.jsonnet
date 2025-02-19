local go_runtime(version, arch) = {
  type: 'pod',
  arch: arch,
  containers: [
    { image: 'golang:' + version },
  ],
};

local task_build_go(setup) = {
  name: 'build go ' + setup.goos,
  runtime: go_runtime('1.19', 'amd64'),
  environment: {
    GOOS: setup.goos,
    BIN: setup.bin,
  },
  steps: [
    { type: 'clone' },
    { type: 'restore_cache', keys: ['cache-sum-{{ md5sum "go.sum" }}', 'cache-date-'], dest_dir: '/go/pkg/mod/cache' },
    {
      type: 'run',
      name: 'build',
      command: |||
        if [ -z ${AGOLA_GIT_TAG} ] || [[ ${AGOLA_GIT_TAG} == *-* ]]; then 
          export VERSION=latest
          export BUILD_VERSION=${AGOLA_GIT_COMMITSHA}
        else
          export VERSION=${AGOLA_GIT_TAG}
          export BUILD_VERSION=${AGOLA_GIT_TAG}
        fi

        echo VERSION: ${VERSION}
        echo BUILD_VERSION: ${BUILD_VERSION}

        go build -ldflags="-X github.com/ercole-io/ercole-agent/v2/cmd.version=${BUILD_VERSION}" -o ${BIN}
      |||,
    },
    {
      type: 'save_to_workspace',
      contents: [{
        source_dir: '.',
        dest_dir: '.',
        paths: [
          setup.bin,
          'Makefile',
          'package/**',
          'fetch/**',
          'sql/**',
          'config.json',
          'LICENSE',  // Needed by windows
        ],
      }],
    },
  ],
  depends: ['staticcheck'],
};

local task_pkg_build_deb() = {
  name: 'pkg build deb',
  runtime: {
    type: 'pod',
    arch: 'amd64',
    containers: [
      { image : 'ubuntu:20.04' }
    ]
  },
  working_dir: '/project',
  environment: {
    WORKSPACE: '/project',
  },
  steps: [
    { type: 'restore_workspace', dest_dir: '.' },
    {
      type: 'run',
      name: 'version',
      command: |||
        if [ -z ${AGOLA_GIT_TAG} ] || [[ ${AGOLA_GIT_TAG} == *-* ]]; then 
          if [[ ${AGOLA_GIT_TAG} == *-rc.* ]]; then
            export VERSION=$(echo $AGOLA_GIT_TAG | sed 's/-/_/g')
          else
            export VERSION=latest
          fi
        else
          export VERSION=$(echo $AGOLA_GIT_TAG | sed 's/-/_/g')
        fi
        echo "VERSION: ${VERSION}"
        echo "export VERSION=${VERSION}" > /tmp/variables
      |||,
    },
    { type: 'run', command: '. /tmp/variables && mkdir -p ercole-agent_${VERSION}_x86_64/opt/ercole-agent && cp ercole-agent ercole-agent_${VERSION}_x86_64/opt/ercole-agent'},
    { type: 'run', command: '. /tmp/variables && mkdir -p ercole-agent_${VERSION}_x86_64/opt/ercole-agent/tmp'},
    { type: 'run', command: '. /tmp/variables && cp config.json ercole-agent_${VERSION}_x86_64/opt/ercole-agent/tmp'},
    { type: 'run', command: '. /tmp/variables && cp -r fetch ercole-agent_${VERSION}_x86_64/opt/ercole-agent'},
    { type: 'run', command: '. /tmp/variables && rm -r ercole-agent_${VERSION}_x86_64/opt/ercole-agent/fetch/windows'},
    { type: 'run', command: '. /tmp/variables && cp -r sql ercole-agent_${VERSION}_x86_64/opt/ercole-agent'},
    { type: 'run', command: '. /tmp/variables && mkdir -p ercole-agent_${VERSION}_x86_64/DEBIAN'},
    {
      type: 'run',
      name: 'create control',
      command: |||
        . /tmp/variables && 
         sed -i "s|VERSION|${VERSION}|g" package/deb/control &&
         cp package/deb/control ercole-agent_${VERSION}_x86_64/DEBIAN
      |||,
    },
    { type: 'run', command: '. /tmp/variables && cp package/deb/preinst ercole-agent_${VERSION}_x86_64/DEBIAN' },
    { type: 'run', command: '. /tmp/variables && cp package/deb/postinst ercole-agent_${VERSION}_x86_64/DEBIAN' },
    {
      type: 'run',
      name: 'deb build',
      command: '. /tmp/variables && dpkg-deb --root-owner-group --build --nocheck ercole-agent_${VERSION}_x86_64',
    },
    { type: 'run', command: 'mkdir dist' },
    { type: 'run', command: '. /tmp/variables && cp ercole-agent_${VERSION}_x86_64.deb dist'},
    { type: 'save_to_workspace', contents: [{ source_dir: './dist/', dest_dir: '/dist/', paths: ['**'] }] },
  ],
  depends: ['build go linux'],
};

local task_pkg_build_rhel(setup) = {
  name: 'pkg build ' + setup.dist,
  runtime: {
    type: 'pod',
    arch: 'amd64',
    containers: [
      { image: setup.pkg_build_image },
    ],
  },
  working_dir: '/project',
  environment: {
    WORKSPACE: '/project',
    DIST: setup.dist,
  },
  steps: [
    { type: 'restore_workspace', dest_dir: '.' },
    {
      type: 'run',
      name: 'version',
      command: |||
        if [ -z ${AGOLA_GIT_TAG} ] || [[ ${AGOLA_GIT_TAG} == *-* ]]; then 
          if [[ ${AGOLA_GIT_TAG} == *-rc.* ]]; then
            export VERSION=$(echo $AGOLA_GIT_TAG | sed 's/-/_/g')
          else
            export VERSION=latest
          fi
        else
          export VERSION=${AGOLA_GIT_TAG}
        fi
        echo VERSION: ${VERSION}
        echo "export VERSION=${VERSION}" > /tmp/variables
      |||,
    },
    {
      type: 'run',
      name: 'sed version',
      command: |||
        source /tmp/variables

        sed -i "s|ERCOLE_VERSION|${VERSION}|g" package/rhel9/ercole-agent.spec
        sed -i "s|ERCOLE_VERSION|${VERSION}|g" package/rhel8/ercole-agent.spec
        sed -i "s|ERCOLE_VERSION|${VERSION}|g" package/rhel7/ercole-agent.spec
        sed -i "s|ERCOLE_VERSION|${VERSION}|g" package/rhel6/ercole-agent.spec
      |||,
    },
    { type: 'run', command: 'rpmbuild --quiet -bl package/${DIST}/ercole-agent.spec || echo ok' },
    { type: 'run', command: 'source /tmp/variables && mkdir -p ~/rpmbuild/SOURCES/ercole-agent-${VERSION}' },
    { type: 'run', command: 'source /tmp/variables && cp -r * ~/rpmbuild/SOURCES/ercole-agent-${VERSION}/' },
    { type: 'run', command: 'source /tmp/variables && tar -C ~/rpmbuild/SOURCES -cvzf ~/rpmbuild/SOURCES/ercole-agent-${VERSION}.tar.gz ercole-agent-${VERSION}' },
    { type: 'run', command: 'pwd; ls && rpmbuild -v -bb package/${DIST}/ercole-agent.spec' },
    { type: 'run', command: 'find ~/rpmbuild/' },
    { type: 'run', command: 'mkdir dist' },
    { type: 'run', command: 'ls ~/rpmbuild/RPMS/x86_64/' },
    { type: 'run', command: 'source /tmp/variables && cd ${WORKSPACE} && cp ~/rpmbuild/RPMS/x86_64/ercole-agent-${VERSION}-1*.x86_64.rpm dist/' },
    { type: 'run', command: 'ls ~/rpmbuild/RPMS/x86_64/ercole-*.rpm' },
    { type: 'run', command: 'file ~/rpmbuild/RPMS/x86_64/ercole-*.rpm' },
    { type: 'run', command: 'cp ~/rpmbuild/RPMS/x86_64/ercole-*.rpm ${WORKSPACE}/dist' },
    { type: 'save_to_workspace', contents: [{ source_dir: './dist/', dest_dir: '/dist/', paths: ['**'] }] },
  ],
  depends: ['build go linux'],
};

local task_deploy_repository(dist) = {
  name: 'deploy repository.ercole.io ' + dist,
  runtime: {
    type: 'pod',
    arch: 'amd64',
    containers: [
      { image: 'curlimages/curl' },
    ],
  },
  environment: {
    REPO_USER: { from_variable: 'repo-user' },
    REPO_TOKEN: { from_variable: 'repo-token' },
    REPO_UPLOAD_URL: { from_variable: 'repo-upload-url' },
    REPO_INSTALL_URL: { from_variable: 'repo-install-url' },
  },
  steps: [
    { type: 'restore_workspace', dest_dir: '.' },
    {
      type: 'run',
      name: 'curl',
      command: |||
        cd dist
        for f in *; do
        	URL=$(curl --user "${REPO_USER}" \
            --upload-file $f ${REPO_UPLOAD_URL} --insecure)
        	echo $URL
        	md5sum $f
        	curl -H "X-API-Token: ${REPO_TOKEN}" \
          -H "Content-Type: application/json" --request POST --data "{ \"filename\": \"$f\", \"url\": \"$URL\" }" \
          ${REPO_INSTALL_URL} --insecure
        done
      |||,
    },
  ],
  depends: ['pkg build ' + dist],
  when: {
    tag: '#.*#',
    branch: 'master',
  },
};

local task_upload_asset(dist) = {
 name: 'upload to github.com ' + dist,
  runtime: {
    type: 'pod',
    arch: 'amd64',
    containers: [
      { image: 'curlimages/curl' },
    ],
  },
 environment: {
    GITHUB_USER: { from_variable: 'github-user' },
    GITHUB_TOKEN: { from_variable: 'github-token' },
  },
steps: [
    { type: 'restore_workspace', dest_dir: '.' },
    {
      type: 'run',
      name: 'upload to github',
      command: |||
          cd dist
          GH_REPO="https://api.github.com/repos/${GITHUB_USER}/ercole-agent/releases"
          if [ ${AGOLA_GIT_TAG} ];
            then GH_TAGS="$GH_REPO/tags/$AGOLA_GIT_TAG" ;
          else
            GH_TAGS="$GH_REPO/latest" ; fi
          response=$(curl -sH "Authorization: token ${GITHUB_TOKEN}" $GH_TAGS)
          eval $(echo "$response" | grep -m 1 "id.:" | grep -w id | tr : = | tr -cd '[[:alnum:]]=')
          for filename in *; do
            REPO_ASSET="https://uploads.github.com/repos/${GITHUB_USER}/ercole-agent/releases/$id/assets?name=$(basename $filename)"
            curl -H POST -H "Authorization: token ${GITHUB_TOKEN}" -H "Content-Type: application/octet-stream" --data-binary @"$filename" $REPO_ASSET
            echo $REPO_ASSET
          done
      |||,
    },
  ],
  depends: ['pkg build ' + dist],
  when: {
    tag: '#.*#',
    branch: 'master',
  },
};

{
  runs: [
    {
      name: 'ercole-agent',
      tasks: [
        {
          name: 'test',
          runtime: {
            type: 'pod',
            arch: 'amd64',
            containers: [
              { image: 'golang:1.19' },
            ],
          },
          steps: [
            { type: 'clone' },
            { type: 'restore_cache', keys: ['cache-sum-{{ md5sum "go.sum" }}', 'cache-date-'], dest_dir: '/go/pkg/mod/cache' },

            { type: 'run', name: 'install golangci-lint', command: 'curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.59.1' },
            { type: 'run', name: 'run golangci-lint', command: 'golangci-lint run' },

            { type: 'run', name: '', command: 'go install go.uber.org/mock/mockgen@v0.3.0' },
            { type: 'run', name: '', command: 'go generate -v ./...' },
            { type: 'run', name: '', command: 'go test -race -coverprofile=coverage.txt -covermode=atomic ./...' },

            { type: 'save_cache', key: 'cache-sum-{{ md5sum "go.sum" }}', contents: [{ source_dir: '/go/pkg/mod/cache' }] },
            { type: 'save_cache', key: 'cache-date-{{ year }}-{{ month }}-{{ day }}', contents: [{ source_dir: '/go/pkg/mod/cache' }] },
          ],
        },
      ] + [
        {
          name: 'staticcheck',
          runtime: {
            type: 'pod',
            arch: 'amd64',
            containers: [
              { image: 'golang:1.22' },
            ],
          },
          steps: [
            { type: 'clone' },
            { type: 'run', name: 'install staticcheck', command: 'go install honnef.co/go/tools/cmd/staticcheck@latest' },
            { type: 'run', name: 'run staticcheck', command: 'staticcheck -f=stylish -tests=false ./...' },
          ],
          depends: ['test'],
        },
      ] + [
        task_build_go(setup)
        for setup in [
          { goos: 'linux', bin: 'ercole-agent' },
          { goos: 'windows', bin: 'ercole-agent.exe' },
        ]
      ] + [
        task_pkg_build_rhel(setup)
        for setup in [
          { pkg_build_image: 'amreo/rpmbuild-centos6', dist: 'rhel6', distfamily: 'rhel' },
          { pkg_build_image: 'amreo/rpmbuild-centos7', dist: 'rhel7', distfamily: 'rhel' },
          { pkg_build_image: 'amreo/rpmbuild-centos8', dist: 'rhel8', distfamily: 'rhel' },
          { pkg_build_image: 'fra.ocir.io/fremyxlx6yog/rpmbuildrhel9', dist: 'rhel9', distfamily: 'rhel' },
        ]
      ] + [
        task_pkg_build_deb()
      ] + [
        {
          name: 'pkg build windows',
          runtime: {
            type: 'pod',
            arch: 'amd64',
            containers: [
              { image: 'amreo/nsis' },
            ],
          },
          working_dir: '/project',
          environment: {
            WORKSPACE: '/project',
            DIST: 'win',
          },
          steps: [
            { type: 'restore_workspace', dest_dir: '.' },
            {
              type: 'run',
              name: 'version',
              command: |||
                if [ -z ${AGOLA_GIT_TAG} ] || [[ ${AGOLA_GIT_TAG} == *-* ]]; then 
                  if [[ ${AGOLA_GIT_TAG} == *-rc.* ]]; then
                    export VERSION=$(echo $AGOLA_GIT_TAG | sed 's/-/_/g')
                  else
                    export VERSION=latest
                  fi
                else
                  export VERSION=${AGOLA_GIT_TAG}
                fi
                echo VERSION: ${VERSION}
                echo "export VERSION=${VERSION}" > /tmp/variables
              |||,
            },
            {
              type: 'run',
              name: 'sed version',
              command: 'source /tmp/variables && sed -i "s|ERCOLE_VERSION|${VERSION}|g" package/win/installer.nsi',
            },
            { type: 'run', command: 'mkdir dist' },
            { type: 'run', command: 'makensis package/win/installer.nsi' },
            { type: 'run', command: 'md5sum ercole-agent.exe' },
            { type: 'run', command: 'source /tmp/variables && cp ercole-agent-setup-${VERSION}.exe dist/' },
            { type: 'save_to_workspace', contents: [{ source_dir: './dist/', dest_dir: '/dist/', paths: ['**'] }] },
          ],
          depends: ['build go windows'],
        },
      ] + [
        task_deploy_repository(dist)
        for dist in ['rhel6', 'rhel7', 'rhel8', 'windows', 'deb', 'rhel9']
      ],
    },
  ],
}
