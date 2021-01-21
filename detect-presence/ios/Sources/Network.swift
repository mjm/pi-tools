import Relay
import Foundation
import Combine

private let graphqlURL = URL(string: "http://100.117.39.47:8080/graphql")!

struct RequestPayload: Encodable {
    var query: String
    var operationName: String
    var variables: VariableData
}

class Network: Relay.Network {
    let url: URL

    init(isDevServer: Bool) {
        if isDevServer {
            NSLog("creating dev server network")
            self.url = URL(string: "http://100.117.39.47:8080/graphql")!
        } else {
            NSLog("creating real network")
            self.url = URL(string: "https://homebase.homelab/graphql")!
        }
    }

    func execute(request: RequestParameters, variables: VariableData, cacheConfig: CacheConfig) -> AnyPublisher<Data, Error> {
        var req = URLRequest(url: graphqlURL)
        req.setValue("application/json", forHTTPHeaderField: "Content-Type")
        req.httpMethod = "POST"

        do {
            let payload = RequestPayload(query: request.text!, operationName: request.name, variables: variables)
            req.httpBody = try JSONEncoder().encode(payload)
        } catch {
            return Fail(error: error).eraseToAnyPublisher()
        }

        return URLSession.shared.dataTaskPublisher(for: req)
            .map { $0.data }
            .mapError { $0 as Error }
            .eraseToAnyPublisher()
    }
}
