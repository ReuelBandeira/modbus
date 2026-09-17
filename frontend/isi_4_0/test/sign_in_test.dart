import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:isi_4_0/providers/lang.dart';
import 'package:isi_4_0/viewmodel/login_view_model.dart';
import 'package:isi_4_0/views/login/sign_in.dart';
import 'package:provider/provider.dart';

void main() {
  final binding = TestWidgetsFlutterBinding.ensureInitialized();
  binding.reset();
  testWidgets('Login Test', (WidgetTester tester) async {
    final loginProvider = LoginViewModel();
    final languageProvider = LanguageProvider();

    await tester.pumpWidget(
      MultiProvider(providers: [
        ChangeNotifierProvider<LoginViewModel>.value(value: loginProvider),
        ChangeNotifierProvider<LanguageProvider>.value(value: languageProvider)
      ], child: const MaterialApp(home: SignIn())),
    );
  });
}
