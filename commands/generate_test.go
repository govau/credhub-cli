package commands_test

import (
	"net/http"

	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
	"github.com/pivotal-cf/credhub-cli/commands"
	"fmt"
)

var _ = Describe("Generate", func() {
	Describe("Without parameters", func() {
		It("uses default parameters", func() {
			setupPasswordPostServer("my-password", "potatoes", generateDefaultTypeRequestJson(`{}`, true))

			session := runCommand("generate", "-n", "my-password")
			Eventually(session).Should(Exit(0))
		})

		It("prints the generated password secret", func() {
			setupPasswordPostServer("my-password", "potatoes", generateDefaultTypeRequestJson(`{}`, true))

			session := runCommand("generate", "-n", "my-password")

			Eventually(session).Should(Exit(0))
			Expect(session.Out).To(Say(responseMyPasswordPotatoes))
		})

		It("can print the generated password secret as JSON", func() {
			setupPasswordPostServer("my-password", "potatoes", generateDefaultTypeRequestJson(`{}`, true))

			session := runCommand("generate", "-n", "my-password", "--output-json")

			Eventually(session).Should(Exit(0))
			Expect(session.Out.Contents()).To(MatchJSON(`{
				"type": "password",
				"updated_at": "` + TIMESTAMP + `",
				"value": "potatoes"
			}`))
		})
	})

	Describe("with a variety of password parameters", func() {
		It("prints the secret", func() {
			setupPasswordPostServer("my-password", "potatoes", generateDefaultTypeRequestJson(`{}`, true))

			session := runCommand("generate", "-n", "my-password", "-t", "password")

			Eventually(session).Should(Exit(0))
			Expect(session.Out).To(Say(responseMyPasswordPotatoes))
		})

		It("can print the secret as JSON", func() {
			setupPasswordPostServer("my-password", "potatoes", generateDefaultTypeRequestJson(`{}`, true))

			session := runCommand(
				"generate",
				"-n", "my-password",
				"-t", "password",
				"--output-json",
			)

			Eventually(session).Should(Exit(0))
			Expect(string(session.Out.Contents())).To(MatchJSON(`{
				"type": "password",
				"updated_at": "` + TIMESTAMP + `",
				"value": "potatoes"
			}`))
		})

		It("with with no-overwrite", func() {
			setupPasswordPostServer("my-password", "potatoes", generateRequestJson("password", `{}`, false))
			session := runCommand("generate", "-n", "my-password", "-t", "password", "--no-overwrite")
			Eventually(session).Should(Exit(0))
		})

		It("including length", func() {
			setupPasswordPostServer("my-password", "potatoes", generateRequestJson("password", `{"length":42}`, true))
			session := runCommand("generate", "-n", "my-password", "-t", "password", "-l", "42")
			Eventually(session).Should(Exit(0))
		})

		It("excluding upper case", func() {
			setupPasswordPostServer("my-password", "potatoes", generateRequestJson("password", `{"exclude_upper":true}`, true))
			session := runCommand("generate", "-n", "my-password", "-t", "password", "--exclude-upper")
			Eventually(session).Should(Exit(0))
		})

		It("excluding lower case", func() {
			setupPasswordPostServer("my-password", "potatoes", generateRequestJson("password", `{"exclude_lower":true}`, true))
			session := runCommand("generate", "-n", "my-password", "-t", "password", "--exclude-lower")
			Eventually(session).Should(Exit(0))
		})

		It("excluding special characters", func() {
			setupPasswordPostServer("my-password", "potatoes", generateRequestJson("password", `{"exclude_special":true}`, true))
			session := runCommand("generate", "-n", "my-password", "-t", "password", "--exclude-special")
			Eventually(session).Should(Exit(0))
		})

		It("excluding numbers", func() {
			setupPasswordPostServer("my-password", "potatoes", generateRequestJson("password", `{"exclude_number":true}`, true))
			session := runCommand("generate", "-n", "my-password", "-t", "password", "--exclude-number")
			Eventually(session).Should(Exit(0))
		})

		It("including only hex", func() {
			setupPasswordPostServer("my-password", "potatoes", generateRequestJson("password", `{"only_hex":true}`, true))
			session := runCommand("generate", "-n", "my-password", "-t", "password", "--only-hex")
			Eventually(session).Should(Exit(0))
		})
	})

	Describe("with a variety of SSH parameters", func() {
		It("prints the SSH key", func() {
			setupRsaSshPostServer("foo-ssh-key", "ssh", "some-public-key", "some-private-key", generateRequestJson("ssh", `{}`, true))

			session := runCommand("generate", "-n", "foo-ssh-key", "-t", "ssh")

			Eventually(session).Should(Exit(0))
			Expect(session.Out).To(Say(responseMySSHFoo))
		})

		It("can print the SSH key as JSON", func() {
			setupRsaSshPostServer("foo-ssh-key", "ssh", "some-public-key", "fake-private-key", generateRequestJson("ssh", `{}`, true))

			session := runCommand("generate", "-n", "foo-ssh-key", "-t", "ssh", "--output-json")

			Eventually(session).Should(Exit(0))
			Expect(string(session.Out.Contents())).To(MatchJSON(`{
				"type": "ssh",
				"updated_at": "` + TIMESTAMP + `",
				"public_key": "some-public-key",
				"private_key": "fake-private-key"
			}`))
		})

		It("with with no-overwrite", func() {
			setupRsaSshPostServer("my-ssh", "ssh", "some-public-key", "some-private-key", generateRequestJson("ssh", `{}`, false))
			session := runCommand("generate", "-n", "my-ssh", "-t", "ssh", "--no-overwrite")
			Eventually(session).Should(Exit(0))
		})

		It("including length", func() {
			setupRsaSshPostServer("my-ssh", "ssh", "some-public-key", "some-private-key", generateRequestJson("ssh", `{"key_length":3072}`, true))
			session := runCommand("generate", "-n", "my-ssh", "-t", "ssh", "-k", "3072")
			Eventually(session).Should(Exit(0))
		})

		It("including comment", func() {
			expectedRequestJson := generateRequestJson("ssh", `{"ssh_comment":"i am an ssh comment"}`, true)
			setupRsaSshPostServer("my-ssh", "ssh", "some-public-key", "some-private-key", expectedRequestJson)
			session := runCommand("generate", "-n", "my-ssh", "-t", "ssh", "-m", "i am an ssh comment")
			Eventually(session).Should(Exit(0))
		})
	})

	Describe("with a variety of RSA parameters", func() {
		It("prints the RSA key", func() {
			setupRsaSshPostServer("foo-rsa-key", "rsa", "some-public-key", "some-private-key", generateRequestJson("rsa", `{}`, true))

			session := runCommand("generate", "-n", "foo-rsa-key", "-t", "rsa")

			Eventually(session).Should(Exit(0))
			Expect(session.Out).To(Say(responseMyRSAFoo))
		})

		It("can print the RSA key as JSON", func() {
			setupRsaSshPostServer("foo-rsa-key", "rsa", "some-public-key", "fake-private-key", generateRequestJson("rsa", `{}`, true))

			session := runCommand("generate", "-n", "foo-rsa-key", "-t", "rsa", "--output-json")

			Eventually(session).Should(Exit(0))
			Expect(string(session.Out.Contents())).To(MatchJSON(`{
				"type": "rsa",
				"updated_at": "` + TIMESTAMP + `",
				"public_key": "some-public-key",
				"private_key": "fake-private-key"
			}`))
		})

		It("with with no-overwrite", func() {
			setupRsaSshPostServer("my-rsa", "rsa", "some-public-key", "some-private-key", generateRequestJson("rsa", `{}`, false))
			session := runCommand("generate", "-n", "my-rsa", "-t", "rsa", "--no-overwrite")
			Eventually(session).Should(Exit(0))
		})

		It("including length", func() {
			setupRsaSshPostServer("my-rsa", "rsa", "some-public-key", "some-private-key", generateRequestJson("rsa", `{"key_length":3072}`, true))
			session := runCommand("generate", "-n", "my-rsa", "-t", "rsa", "-k", "3072")
			Eventually(session).Should(Exit(0))
		})
	})

	Describe("with a variety of certificate parameters", func() {
		It("prints the certificate", func() {
			expectedRequestJson := generateRequestJson("certificate", `{"common_name":"common.name.io"}`, true)
			setupCertificatePostServer("my-secret", "my-ca", "my-cert", "my-priv", expectedRequestJson)

			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--common-name", "common.name.io")

			Eventually(session).Should(Exit(0))
			Expect(session.Out).To(Say(responseMyCertificate))
		})

		It("can print the certificate as JSON", func() {
			expectedRequestJson := generateRequestJson("certificate", `{"common_name":"common.name.io"}`, true)
			setupCertificatePostServer("my-secret", "my-ca", "my-cert", "my-priv", expectedRequestJson)

			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--common-name", "common.name.io", "--output-json")

			Eventually(session).Should(Exit(0))
			Expect(string(session.Out.Contents())).To(MatchJSON(`{
				"type": "certificate",
				"updated_at": "` + TIMESTAMP + `",
				"ca": "my-ca",
				"certificate": "my-cert",
				"private_key": "my-priv"
			}`))
		})

		It("including common name", func() {
			expectedRequestJson := generateRequestJson("certificate", `{"common_name":"common.name.io"}`, true)
			setupCertificatePostServer("my-secret", "potatoes-ca", "potatoes-cert", "potatoes-priv-key", expectedRequestJson)
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--common-name", "common.name.io")
			Eventually(session).Should(Exit(0))
		})

		It("including common name with no-overwrite", func() {
			expectedRequestJson := generateRequestJson("certificate", `{"common_name":"common.name.io"}`, false)
			setupCertificatePostServer("my-secret", "potatoes-ca", "potatoes-cert", "potatoes-priv-key", expectedRequestJson)
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--common-name", "common.name.io", "--no-overwrite")
			Eventually(session).Should(Exit(0))
		})

		It("including organization", func() {
			expectedRequestJson := generateRequestJson("certificate", `{"organization":"organization.io"}`, true)
			setupCertificatePostServer("my-secret", "potatoes-ca", "potatoes-cert", "potatoes-priv-key", expectedRequestJson)
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--organization", "organization.io")
			Eventually(session).Should(Exit(0))
		})

		It("including organization unit", func() {
			expectedRequestJson := generateRequestJson("certificate", `{"organization_unit":"My Unit"}`, true)
			setupCertificatePostServer("my-secret", "potatoes-ca", "potatoes-cert", "potatoes-priv-key", expectedRequestJson)
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--organization-unit", "My Unit")
			Eventually(session).Should(Exit(0))
		})

		It("including locality", func() {
			expectedRequestJson := generateRequestJson("certificate", `{"locality":"My Locality"}`, true)
			setupCertificatePostServer("my-secret", "potatoes-ca", "potatoes-cert", "potatoes-priv-key", expectedRequestJson)
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--locality", "My Locality")
			Eventually(session).Should(Exit(0))
		})

		It("including state", func() {
			expectedRequestJson := generateRequestJson("certificate", `{"state":"My State"}`, true)
			setupCertificatePostServer("my-secret", "potatoes-ca", "potatoes-cert", "potatoes-priv-key", expectedRequestJson)
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--state", "My State")
			Eventually(session).Should(Exit(0))
		})

		It("including country", func() {
			expectedRequestJson := generateRequestJson("certificate", `{"country":"My Country"}`, true)
			setupCertificatePostServer("my-secret", "potatoes-ca", "potatoes-cert", "potatoes-priv-key", expectedRequestJson)
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--country", "My Country")
			Eventually(session).Should(Exit(0))
		})

		It("including multiple alternative names", func() {
			expectedRequestJson := generateRequestJson("certificate", `{"alternative_names": [ "Alt1", "Alt2" ]}`, true)
			setupCertificatePostServer("my-secret", "potatoes-ca", "potatoes-cert", "potatoes-priv-key", expectedRequestJson)
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--alternative-name", "Alt1", "--alternative-name", "Alt2")
			Eventually(session).Should(Exit(0))
		})

		It("including key length", func() {
			expectedRequestJson := generateRequestJson("certificate", `{"key_length":2048}`, true)
			setupCertificatePostServer("my-secret", "potatoes-ca", "potatoes-cert", "potatoes-priv-key", expectedRequestJson)
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--key-length", "2048")
			Eventually(session).Should(Exit(0))
		})

		It("including duration", func() {
			expectedRequestJson := generateRequestJson("certificate", `{"duration":1000}`, true)
			setupCertificatePostServer("my-secret", "potatoes-ca", "potatoes-cert", "potatoes-priv-key", expectedRequestJson)
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--duration", "1000")
			Eventually(session).Should(Exit(0))
		})

		It("including certificate authority", func() {
			expectedRequestJson := generateRequestJson("certificate", `{"ca":"my_ca"}`, true)
			setupCertificatePostServer("my-secret", "potatoes-ca", "potatoes-cert", "potatoes-priv-key", expectedRequestJson)
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--ca", "my_ca")
			Eventually(session).Should(Exit(0))
		})
	})

	Describe("Help", func() {
		ItBehavesLikeHelp("generate", "n", func(session *Session) {
			Expect(session.Err).To(Say("generate"))
			Expect(session.Err).To(Say("name"))
			Expect(session.Err).To(Say("length"))
		})

		It("short flags", func() {
			Expect(commands.GenerateCommand{}).To(SatisfyAll(
				commands.HaveFlag("name", "n"),
				commands.HaveFlag("type", "t"),
				commands.HaveFlag("no-overwrite", "O"),
				commands.HaveFlag("length", "l"),
				commands.HaveFlag("exclude-special", "S"),
				commands.HaveFlag("exclude-number", "N"),
				commands.HaveFlag("exclude-upper", "U"),
				commands.HaveFlag("exclude-lower", "L"),
				commands.HaveFlag("only-hex", "H"),
				commands.HaveFlag("common-name", "c"),
				commands.HaveFlag("organization", "o"),
				commands.HaveFlag("organization-unit", "u"),
				commands.HaveFlag("locality", "i"),
				commands.HaveFlag("state", "s"),
				commands.HaveFlag("country", "y"),
				commands.HaveFlag("alternative-name", "a"),
				commands.HaveFlag("key-length", "k"),
				commands.HaveFlag("duration", "d"),
			))
		})

		It("displays missing 'n' option as required parameters", func() {
			session := runCommand("generate")

			Eventually(session).Should(Exit(1))

			if runtime.GOOS == "windows" {
				Expect(session.Err).To(Say("the required flag `/n, /name' was not specified"))
			} else {
				Expect(session.Err).To(Say("the required flag `-n, --name' was not specified"))
			}
		})

		It("displays the server provided error when an error is received", func() {
			server.AppendHandlers(
				RespondWith(http.StatusBadRequest, `{"error": "you fail."}`),
			)

			session := runCommand("generate", "-n", "my-value")

			Eventually(session).Should(Exit(1))

			Expect(session.Err).To(Say("you fail."))
		})
	})
})

func setupPasswordPostServer(name string, value string, requestJson string) {
	server.AppendHandlers(
		CombineHandlers(
			VerifyRequest("POST", fmt.Sprintf("/api/v1/data/%s", name)),
			VerifyJSON(requestJson),
			RespondWith(http.StatusOK, fmt.Sprintf(STRING_SECRET_RESPONSE_JSON, "password", name, value)),
		),
	)
}

func setupRsaSshPostServer(name string, contentType string, publicKey string, privateKey string, requestJson string) {
	server.AppendHandlers(
		CombineHandlers(
			VerifyRequest("POST", fmt.Sprintf("/api/v1/data/%s", name)),
			VerifyJSON(requestJson),
			RespondWith(http.StatusOK, fmt.Sprintf(RSA_SSH_SECRET_RESPONSE_JSON, contentType, name, publicKey, privateKey)),
		),
	)
}

func setupCertificatePostServer(name string, ca string, certificate string, privateKey string, requestJson string) {
	server.AppendHandlers(
		CombineHandlers(
			VerifyRequest("POST", fmt.Sprintf("/api/v1/data/%s", name)),
			VerifyJSON(requestJson),
			RespondWith(http.StatusOK, fmt.Sprintf(CERTIFICATE_SECRET_RESPONSE_JSON, name, ca, certificate, privateKey)),
		),
	)
}

func generateRequestJson(secretType string, params string, overwrite bool) string {
	return fmt.Sprintf(GENERATE_SECRET_REQUEST_JSON, secretType, overwrite, params)
}

func generateDefaultTypeRequestJson(params string, overwrite bool) string {
	return fmt.Sprintf(GENERATE_DEFAULT_TYPE_REQUEST_JSON, overwrite, params)
}
