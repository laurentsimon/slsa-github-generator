package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	/// Builder information.
	builderReleaseBinary string
	builderRepository    string
	builderRef           string
	// Verifier information.
	verifierRespository        string
	verifierRealeaseBinary     string
	verifierRealseBinarySHA256 string
	verifierReleaseTag         string
)

var errorInvalidRef = errors.New("invalid ref")

func usage(p string) {
	fmt.Printf(`Usage:
	 %s --builder-binary <> --builder-repository <> --builder-ref <>
	 	--verifier-repository <> --verifier-binary <> --verifier-binary-sha256<> --verifier-tag <>
		 `,
		p)
	os.Exit(-1)
}

func check(e error) {
	if e != nil {
		fmt.Fprintf(os.Stderr, e.Error())
		os.Exit(1)
	}
}

func main() {
	// Usage: go run . --builder-binary=slsa-builder-go-linux-amd64 --builder-repository=laurentsimon/slsa-github-generator --builder-ref="refs/tags/v0.0.21" --verifier-repository=laurentsimon/slsa-verifier --verifier-binary=slsa-verifier-linux-amd64 --verifier-binary-sha256=fb743bc6bb56908d590da66bfe5c266d003aa226b30fcada5f7b9e4aea43b52b --verifier-tag=v0.0.4

	// Builder.
	flag.StringVar(&builderReleaseBinary, "builder-binary", "", "Name of the builder's binary in release assets")
	flag.StringVar(&builderRepository, "builder-repository", "", "The builder's repository.")
	flag.StringVar(&builderRef, "builder-ref", "", "The builder's ref.")
	// Verifier.
	flag.StringVar(&verifierRespository, "verifier-repository", "", "The verifier's repository")
	flag.StringVar(&verifierRealeaseBinary, "verifier-binary", "", "Name of the verifier's binary in release assets")
	flag.StringVar(&verifierRealseBinarySHA256, "verifier-binary-sha256", "", "SHA256 of the verifier's binary.")
	flag.StringVar(&verifierReleaseTag, "verifier-tag", "", "Verifier's release tag (e.g., v1.2.3).")

	flag.Parse()

	if builderReleaseBinary == "" || builderRepository == "" || builderRef == "" ||
		verifierRespository == "" || verifierRealeaseBinary == "" || verifierRealseBinarySHA256 == "" ||
		verifierReleaseTag == "" {
		usage(os.Args[0])
	}

	fmt.Printf("Running with arguments:\n")
	fmt.Printf("builder-binary: %s\n", builderReleaseBinary)
	fmt.Printf("builder-repository: %s\n", builderRepository)
	fmt.Printf("builder-ref: %s\n", builderRef)
	fmt.Printf("verifier-repository: %s\n", verifierRespository)
	fmt.Printf("verifier-binary: %s\n", verifierRealeaseBinary)
	fmt.Printf("verifier-binary-sha256: %s\n", verifierRealseBinarySHA256)
	fmt.Printf("verifier-tag: %s\n", verifierReleaseTag)
	fmt.Println()

	// Verify the builder's ref. It must be a version
	builderTag, err := extractBuilderRef(builderRef)
	check(err)

	fmt.Println("Builder tag:", builderTag)
	/*
		resp, err := http.Get("http://example.com/")
		if err != nil {
			// handle error
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)*/
}

func extractBuilderRef(ref string) (string, error) {
	if !strings.HasPrefix(ref, "refs/tags/") {
		return "", fmt.Errorf("%w: %s is not a valid tag", errorInvalidRef, ref)
	}

	tag := strings.TrimPrefix(ref, "refs/tags/")
	regex := regexp.MustCompile(`^v\d*(\.([\d]{1,})){0,2}$`)
	match := regex.MatchString(tag)
	if !match {
		return "", fmt.Errorf("%w: %s is not a of the form vX.Y.Z", errorInvalidRef, tag)
	}
	return tag, nil
}

/*


echo "Builder version: $BUILDER_TAG"

echo "BUILDER_REPOSITORY: $BUILDER_REPOSITORY"

# Fetch the release binary and provenance.
gh release -R "$BUILDER_REPOSITORY" download "$BUILDER_TAG" -p "$BUILDER_RELEASE_BINARY*" || exit 10

# Fetch the verifier at the right hash.
gh release -R "$VERIFIER_REPOSITORY" download "$VERIFIER_RELEASE" -p "$VERIFIER_RELEASE_BINARY" || exit 11
COMPUTED_HASH=$(sha256sum "$VERIFIER_RELEASE_BINARY" | awk '{print $1}')
echo "verifier hash computed is $COMPUTED_HASH"
echo "$VERIFIER_RELEASE_BINARY_SHA256 $VERIFIER_RELEASE_BINARY" | sha256sum --strict --check --status || exit 4
echo "verifier hash verification has passed"

# Verify the provenance of the builder.
chmod a+x "$VERIFIER_RELEASE_BINARY"
./"$VERIFIER_RELEASE_BINARY" --branch "main" \
                            --tag "$BUILDER_TAG" \
                            --artifact-path "$BUILDER_RELEASE_BINARY" \
                            --provenance "$BUILDER_RELEASE_BINARY.intoto.jsonl" \
                            --source "github.com/$BUILDER_REPOSITORY" || exit 6

BUILDER_COMMIT=$(gh api /repos/"$BUILDER_REPOSITORY"/git/ref/tags/"$BUILDER_TAG" | jq -r '.object.sha')
PROVENANCE_COMMIT=$(cat "$BUILDER_RELEASE_BINARY.intoto.jsonl" | jq -r '.payload' | base64 -d | jq -r '.predicate.materials[0].digest.sha1')
if [[ "$BUILDER_COMMIT" != "$PROVENANCE_COMMIT" ]]; then
    echo "Builder commit sha $BUILDER_COMMIT != provenance material $PROVENANCE_COMMIT"
    exit 5
fi

#TODO: verify the command
echo "Builder provenance verified at tag $BUILDER_TAG and commit $BUILDER_COMMIT"

*/
