echo "" > coverage.txt

# if [ -f auth_coverage.out ]; then
#   cat auth_coverage.out >> coverage.txt
#   rm auth_coverage.out
# fi

# Locally precovered strategies
strategies=( addin adfs fba ntml saml tmg )
for strategy in "${strategies[@]}"
do
  auth_coverage_file="auth/${strategy}/coverage.data"
  if [ -f auth_coverage_file ]; then
    cat auth_coverage_file >> coverage.txt
  fi
done

if [ -f api_coverage.out ]; then
  cat api_coverage.out >> coverage.txt
  rm api_coverage.out
fi