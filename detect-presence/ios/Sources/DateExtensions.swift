import Foundation

private let dateFormatter = ISO8601DateFormatter()

extension String {
    var asDate: Date? {
        dateFormatter.date(from: self)
    }
}

extension Date {
    var asString: String {
        dateFormatter.string(from: self)
    }
}
