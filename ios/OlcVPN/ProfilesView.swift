import SwiftUI
import VPNCore

struct ProfilesView: View {
    @State private var profiles: [ProfileItem] = []
    @State private var showingAddProfile = false

    var body: some View {
        NavigationView {
            List {
                ForEach(profiles) { profile in
                    ProfileRow(profile: profile)
                        .swipeActions {
                            Button(role: .destructive) {
                                deleteProfile(profile)
                            } label: {
                                Label("Delete", systemImage: "trash")
                            }
                        }
                }
            }
            .navigationTitle("Profiles")
            .toolbar {
                Button(action: { showingAddProfile = true }) {
                    Image(systemName: "plus")
                }
            }
            .sheet(isPresented: $showingAddProfile) {
                AddProfileView(profiles: $profiles)
            }
            .onAppear {
                loadProfiles()
            }
            .background(Color(red: 0.05, green: 0.05, blue: 0.05).ignoresSafeArea())
        }
    }

    private func loadProfiles() {
        // TODO: Load from VPNCore
        profiles = []
    }

    private func deleteProfile(_ profile: ProfileItem) {
        // TODO: Delete via VPNCore
        profiles.removeAll { $0.id == profile.id }
    }
}

struct ProfileRow: View {
    let profile: ProfileItem

    var body: some View {
        VStack(alignment: .leading, spacing: 4) {
            Text(profile.name)
                .font(.headline)
                .foregroundColor(.white)

            HStack {
                Text(profile.engine.uppercased())
                    .font(.caption)
                    .padding(.horizontal, 8)
                    .padding(.vertical, 2)
                    .background(engineColor)
                    .foregroundColor(.black)
                    .cornerRadius(4)

                Text(profile.protocol)
                    .font(.caption)
                    .foregroundColor(.gray)
            }
        }
        .padding(.vertical, 4)
        .listRowBackground(Color(red: 0.1, green: 0.1, blue: 0.1))
    }

    private var engineColor: Color {
        profile.engine == "singbox" ? Color(red: 0, green: 0.898, blue: 1) : Color(red: 0, green: 0.784, blue: 0.325)
    }
}

struct AddProfileView: View {
    @Environment(\.dismiss) var dismiss
    @Binding var profiles: [ProfileItem]
    @State private var uri = ""
    @State private var errorMessage = ""

    var body: some View {
        NavigationView {
            Form {
                Section(header: Text("Import from URI")) {
                    TextField("vless://... or olcrtc://...", text: $uri)
                        .autocapitalization(.none)
                        .disableAutocorrection(true)
                }

                if !errorMessage.isEmpty {
                    Section {
                        Text(errorMessage)
                            .foregroundColor(.red)
                    }
                }

                Section {
                    Button("Import") {
                        importProfile()
                    }
                    .disabled(uri.isEmpty)
                }
            }
            .navigationTitle("Add Profile")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .cancellationAction) {
                    Button("Cancel") {
                        dismiss()
                    }
                }
            }
        }
    }

    private func importProfile() {
        // TODO: Import via VPNCore
        errorMessage = ""
        dismiss()
    }
}

struct ProfileItem: Identifiable {
    let id: String
    let name: String
    let engine: String
    let `protocol`: String
}
