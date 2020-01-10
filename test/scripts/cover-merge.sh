echo "mode: atomic" > coverage.txt

# Locally precovered strategies
strategies=( addin adfs anon fba ntml saml tmg )
for strategy in "${strategies[@]}"
do
  auth_coverage_file="auth/${strategy}/coverage.data"
  echo $auth_coverage_file
  if [ -f $auth_coverage_file ]; then
    cat $auth_coverage_file \
      | egrep -v '^mode.*' \
      >> coverage.txt
  fi
done

if [ -f auth_coverage.out ]; then
  cat auth_coverage.out \
    | egrep -v '^mode.*' \
    | egrep -v '^github.com/koltyakov/gosip/auth/.*' \
    | egrep -v '^github.com/koltyakov/gosip/api/.*' \
    >> coverage.txt
  rm auth_coverage.out
fi

if [ -f api_coverage.out ]; then
  cat api_coverage.out \
    | egrep -v '^mode.*' \
    >> coverage.txt
  rm api_coverage.out
fi