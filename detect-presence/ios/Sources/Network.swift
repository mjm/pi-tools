import Relay
import Foundation
import Combine

let myRelayEnvironment = Environment(
    network: MyNetwork(),
    store: Store()
)

private let graphqlURL = URL(string: "http://100.117.39.47:8080/graphql")!

struct RequestPayload: Encodable {
    var query: String
    var operationName: String
    var variables: VariableData
}

class MyNetwork: Network {
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
