#!/usr/bin/env bash
# Smoke tests against a running back end on a freshly seeded database.
# Usage: back-end/scripts/smoke.sh [base url, default http://localhost:5000/api/1]

set -u
BASE="${1:-http://localhost:5000/api/1}"

MONACO=c39a7593-36c5-4ceb-8d07-6bc82f808381
VER=b11a43a8-4e36-47a8-8e68-301fc1cae6ef
BAG=a22b54b9-5f47-48b9-9f79-412fd2dbf7f0
ARTICLE_MONACO=75000000-0000-4000-8000-000000000001
ARTICLE_DRAFT=75000000-0000-4000-8000-000000000003
ARTICLE_FUTURE=75000000-0000-4000-8000-000000000006
UNKNOWN=00000000-0000-4000-8000-000000000099

pass=0
fail=0

# check <expected status> <path> [substring the body must contain] [substring the body must NOT contain]
check() {
  local want_status=$1 path=$2 want=${3:-} unwanted=${4:-}
  local body status
  body=$(curl -s -w '\n%{http_code}' "$BASE$path")
  status=${body##*$'\n'}
  body=${body%$'\n'*}

  local problems=()
  [[ $status == "$want_status" ]] || problems+=("status $status, want $want_status")
  [[ -z $want || $body == *"$want"* ]] || problems+=("body lacks: $want")
  [[ -z $unwanted || $body != *"$unwanted"* ]] || problems+=("body has: $unwanted")

  if ((${#problems[@]} == 0)); then
    pass=$((pass + 1))
    echo "ok    $path"
  else
    fail=$((fail + 1))
    echo "FAIL  $path -> $(IFS='; '; echo "${problems[*]}")"
    echo "      ${body:0:300}"
  fi
}

# 1-3: lists without params, filters, bad input
check 200 "/races" '"name":"24 Hours of Le Mans"'
check 200 "/races?filter=category_id,eq,f1" '"name":"Monaco Grand Prix"' 'Le Mans'
check 200 "/races?filter=season,gte,2026" 'Monaco'
check 400 "/races?filter=season,eq,abc" '"status":400'
check 400 "/races?filter=foo,eq,1" 'invalid filter column'
check 400 "/races?page=0" 'invalid page'
check 400 "/races?sort=nope,asc" 'invalid sort column'
check 200 "/categories?filter=is_active,eq,true" '"id":"f1"'
check 400 "/categories?filter=is_active,eq,maybe"
check 200 "/drivers?filter=last_name,ilike,ver" '"code":"VER"'

# 4-5: get by id
check 200 "/races/$MONACO" '"lap_number":52'
check 400 "/races/not-a-uuid" 'not a UUID'
check 404 "/races/$UNKNOWN" '"status":404'
check 200 "/categories/f1" '"name":"Formula 1"'
check 200 "/categories/F1" '"id":"f1"'
check 404 "/categories/xyz"
check 200 "/cars/74000000-0000-4000-8000-000000000001" '"tyre_supplier":"Pirelli"'
check 404 "/nope" 'resource not found'

# 6-7: articles hide drafts and future-dated ones, texts by article
check 200 "/articles" "$ARTICLE_MONACO" "$ARTICLE_DRAFT"
check 200 "/articles" '' "$ARTICLE_FUTURE"
check 200 "/articles?sort=published_at,desc" "$ARTICLE_MONACO"
check 404 "/articles/$ARTICLE_DRAFT"
check 404 "/articles/$ARTICLE_FUTURE"
check 200 "/article-texts?filter=article_id,eq,$ARTICLE_MONACO" '"language":"pt-br"'
check 200 "/article-texts?filter=article_id,eq,$ARTICLE_DRAFT" '[]'
check 200 "/article-texts?filter=title,ilike,draft" '[]'

# 8: pagination and related lookups by id
check 200 "/drivers?page=2" '"last_name":"Verstappen"'
check 200 "/drivers?page=2" '' '"last_name":"Acosta"'
check 200 "/drivers?filter=id,in,$VER|$BAG" '"code":"BAG"'
check 400 "/drivers?filter=id,in,$VER|nope"

# 9: article search
check 200 "/articles/search?q=pole&lang=en" '{"text":"pole","match":true}'
check 200 "/articles/search?q=verstappen&lang=pt-BR" '"language":"pt-br"'
check 200 "/articles/search?q=verstappen&lang=en" '"title":"Verstappen masters Monaco"' '"language":"pt-br"'
check 200 "/articles/search?q=%22single+pit+stop%22" "$ARTICLE_MONACO"
check 200 "/articles/search?q=draft+or+scheduled" '[]'
check 200 "/articles/search?q=controlled" '' '**'
check 200 "/articles/search?lang=pt-br" '"title":"A hist'
check 200 "/articles/search" '' "$ARTICLE_FUTURE"
check 400 "/articles/search?lang=not_a_lang" 'invalid lang'
check 400 "/articles/search?page=0" 'invalid page'
check 400 "/articles/search?q=$(printf 'a%.0s' {1..201})" 'invalid q'

echo
echo "$pass passed, $fail failed"
((fail == 0))
