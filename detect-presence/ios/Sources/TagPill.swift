import SwiftUI

struct TagPill: View {
    let text: String

    var body: some View {
        HStack(spacing: 6) {
            Image(systemName: "tag.fill")
            Text(verbatim: text)
        }
        .font(Font.caption.weight(.medium))
        .foregroundColor(Color("TagForeground"))
        .padding(.vertical, 4)
        .padding(.horizontal, 8)
        .background(Color("TagBackground").cornerRadius(8))
    }
}
