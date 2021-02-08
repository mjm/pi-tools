import SwiftUI
import Combine
import AuthenticationServices

struct LoginView: View {
    @ObservedObject var model: AppModel

    var body: some View {
        VStack(spacing: 16) {
            Image(systemName: "person.circle")
                .resizable()
                .frame(width: 50, height: 50)
                .foregroundColor(.blue)

            VStack(spacing: 8) {
                Text("You must be logged in to use this app.")
                    .foregroundColor(.secondary)
                    .font(.title3)
                    .multilineTextAlignment(.center)
                    .padding()
                Button {
                    model.authenticator.logIn()
                } label: {
                    Text("Log In")
                        .foregroundColor(.white)
                        .padding()
                        .background(Color.blue)
                        .clipShape(RoundedRectangle(cornerRadius: 8))
                }
            }
        }
    }
}
