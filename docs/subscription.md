# MangaPlus subscription

MangaPlus has no username/password login. Every install of the official app registers itself as a
device and gets a **device secret** (32 hex characters), which is sent with every API request. A
subscription bought in the app is tied to that device secret.

Without a secret, mangoplus registers a new device on the first request, the same way the app does on
a fresh install. That device is on the free (`basic`) plan. To get your subscription's chapters, give
mangoplus the secret of the app install where your subscription is active.

Adapted from the guide in [frank_mangaplus](https://github.com/akitaonrails/frank_mangaplus)
(`docs/android-secret.md`) by Fabio Akita, MIT licensed. See [Credits](../README.md#credits).

## Getting the secret

The app keeps the secret in its private storage, which can only be read with root. Either use a rooted
phone, or a rooted Android emulator signed in to the same Google account as the subscription.

Where it's stored depends on the app version:

- **Current app** (checked 2026-10): a DataStore file, `files/settings`. It's binary, not text, so
  printing it shows some junk around the values. Look for the 32 hex characters right after `secret`.
- **Older app versions**: `shared_prefs/config.xml`, in `<string name="secret">...</string>`.

### Rooted phone

```sh
# current app
adb shell "su -c 'cat /data/data/jp.co.shueisha.mangaplus/files/settings'"
# older app versions
adb shell "su -c 'cat /data/data/jp.co.shueisha.mangaplus/shared_prefs/config.xml'"
```

Keep the same install afterwards: reinstalling the app or clearing its data will probably create a
new secret.

### Rooted emulator

Takes about 20 minutes the first time. Once you have the secret the emulator can be deleted.

1. Install [Android Studio](https://developer.android.com/studio). In the SDK Manager install an
   Android SDK Platform (33 or newer), the **Google Play** `x86_64` system image (the one with Play
   Store) and the Android Emulator.
2. Create a virtual device (e.g. Pixel 6) with that Google Play image, boot it and sign in to Google
   Play with the account that has the subscription. Don't install MangaPlus yet.
3. Root it with [rootAVD](https://gitlab.com/newbit/rootAVD), which patches the emulator with Magisk:

   ```sh
   git clone https://gitlab.com/newbit/rootAVD.git
   cd rootAVD
   ./rootAVD.sh ListAllAVDs
   # pick the ramdisk path of your device from the output
   ./rootAVD.sh <path-to-ramdisk.img>
   ```

   Restart the emulator, the Magisk app should say "Installed". Check with
   `adb shell su -c "id"`, it should print `uid=0(root)`.
4. Install **MANGA Plus by SHUEISHA** from the Play Store, open it, sign in, restore the
   subscription and read a chapter so the app is fully set up.
5. Read the secret, same as with a rooted phone:

   ```sh
   adb shell "su -c 'cat /data/data/jp.co.shueisha.mangaplus/files/settings'"
   ```

   Keep the emulator's install if you want to read the secret again later.

## Using it

- **mangoplus**: set `Options.Secret` and keep `Options.Viewer = ViewerApp` (the default). Only the app
  viewer sets the cookie that subscription-only images need.
- **mangal**: save the secret alone in a file only you can read, and point
  `providers.mangaplus.secret_file` to it. If the file is missing mangal logs a warning and uses the
  free tier. Keep `providers.mangaplus.viewer = 'app'` (the default).

  ```sh
  (umask 077 && echo '<secret>' > ~/.config/mangal/mangaplus-secret)
  ```

  ```toml
  [providers.mangaplus]
  secret_file = '~/.config/mangal/mangaplus-secret'
  ```

`User.Subscription()` shows the plan of the secret in use: `basic`, `standard` or `deluxe`.

`Chapter.ChapterType` depends on who asks: for a subscribed account the API marks chapters as `FREE`,
for a free device the same chapters can be `STANDARD` or `DELUXE` (e.g. Dragon Ball Super). Without
access the pages request fails with "Invalid user access(11301)".

## When it stops working

- Using "Restore Subscription" in the app (on any device) invalidates the old secret, get it again.
- A lapsed renewal downgrades the account to `basic` without any error, subscription chapters just
  start failing. Check the plan with `User.Subscription()`.
- The API locks a device out for 10+ minutes after ~10 requests in a few seconds ("Invalid
  Parameter (10522)"). mangoplus waits 500ms between requests (`Options.RequestInterval`) to stay
  below that.
