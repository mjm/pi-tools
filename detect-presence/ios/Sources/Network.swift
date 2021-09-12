import Relay
import Foundation
import Combine

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
            self.url = URL(string: "http://mars.home.mattmoriarity.com:8080/graphql")!
        } else {
            NSLog("creating real network")
            self.url = URL(string: "https://homelab.home.mattmoriarity.com/graphql")!
        }
    }

    func execute(request: RequestParameters, variables: VariableData, cacheConfig: CacheConfig) -> AnyPublisher<Data, Error> {
        var req = URLRequest(url: self.url)
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
