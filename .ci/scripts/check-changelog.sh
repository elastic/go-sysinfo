#!/usr/bin/env bash
set -uo pipefail

if [[ -z "${PR_NUMBER}" ]]; then
  echo "PR_NUMBER must be set."
  exit 1
fi

if [[ -z "${BASE_REF}" ]]; then
  echo "BASE_REF must be set."
  exit 1
fi

docs_url="https://github.com/GoogleCloudPlatform/magic-modules/blob/2834761fec3acbf35cacbffe100530f82eada650/.ci/RELEASE_NOTES_GUIDE.md#expected-format"

# Version of https://github.com/hashicorp/go-changelog.
go_changelog_version=ba40b3a
go_changelog_check="go run github.com/hashicorp/go-changelog/cmd/changelog-check@${go_changelog_version}"

expected_changelog_file=.changelog/${PR_NUMBER}.txt

# Verify file is present.
if [ ! -e "${expected_changelog_file}" ]; then
  echo "Changelog file missing at ${expected_changelog_file}.

Please add a changelog entry following the format described [here](${docs_url}).

If this change does not require a changelog entry then label the pull request
with skip-changelog.
" >> $GITHUB_STEP_SUMMARY
  exit 1
fi

# Check the format.
if ! ${go_changelog_check} "${expected_changelog_file}"; then
  echo "Changelog format is invalid. See build log." >> $GITHUB_STEP_SUMMARY
  exit 1
fi

echo "Changelog is valid." >> $GITHUB_STEP_SUMMARY
