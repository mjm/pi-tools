import Foundation
import Dispatch

class BackoffExecutor {
    typealias WorkItem = (@escaping (Bool) -> Void) -> Void

    private let queue = DispatchQueue(label: "BackoffExecutor")

    func enqueue(initialDelay: TimeInterval, factor: Double = 2, maxRetries: Int = 5, workItem: @escaping WorkItem) {
        var nextDelay = initialDelay
        var retries = 0

        func performRetry() {
            workItem { succeeded in
                if !succeeded {
                    guard retries < maxRetries else {
                        return
                    }

                    retries += 1
                    let newDelay = nextDelay
                    nextDelay = nextDelay * factor

                    self.queue.asyncAfter(deadline: .now() + newDelay, execute: performRetry)
                }
            }
        }

        queue.async(execute: performRetry)
    }
}
