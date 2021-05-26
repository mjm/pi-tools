import SwiftUI
import Combine
import AuthenticationServices

class Authenticator: NSObject, ObservableObject, ASWebAuthenticationPresentationContextProviding {
    private var cancellables = Set<AnyCancellable>()

    func logIn() {
        Future<URL, Error> { completion in
            let authUrl = URL(string: "https://auth.home.mattmoriarity.com/webauthn/login_app")!

            let authSession = ASWebAuthenticationSession(
                url: authUrl,
                callbackURLScheme: "x-presence-app"
            ) { url, error in
                if let error = error {
                    completion(.failure(error))
                } else if let url = url {
                    completion(.success(url))
                }
            }

            authSession.presentationContextProvider = self
            authSession.prefersEphemeralWebBrowserSession = true
            authSession.start()
        }.sink { completion in
            if case .failure(let error) = completion {
                NSLog("Error signing in: \(error)")
            }
        } receiveValue: { url in
            NSLog("Received URL \(url) from authentication session")
            guard let components = URLComponents(url: url, resolvingAgainstBaseURL: false) else {
                NSLog("couldn't create URL components")
                return
            }

            guard components.path == "/login-succeeded" else {
                NSLog("unexpected path for URL: \(components.path)")
                return
            }

            guard let cookie = components.queryItems?.first(where: { $0.name == "cookie" })?.value else {
                NSLog("couldn't find a URL query item named cookie")
                return
            }

            let cookieParts = cookie.split(separator: "=", maxSplits: 1)
            let newCookie = HTTPCookie(properties: [
                .name: String(cookieParts[0]),
                .value: String(cookieParts[1]),
                .domain: ".home.mattmoriarity.com",
                .path: "/",
                .expires: NSDate().addingTimeInterval(31536000)
            ])!
            HTTPCookieStorage.shared.setCookie(newCookie)
            NSLog("Set \(newCookie.name) cookie in shared cookie storage")

            NSLog("Now the shared cookie storage has: \(HTTPCookieStorage.shared.cookies!)")
        }.store(in: &cancellables)
    }

    func presentationAnchor(for session: ASWebAuthenticationSession) -> ASPresentationAnchor {
        UIWindow()
    }
}
