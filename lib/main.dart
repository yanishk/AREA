import 'package:flutter/material.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';

import "login_page.dart";

Future main() async {
  await dotenv.load(fileName: ".env");
  runApp(MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      debugShowCheckedModeBanner: false,
      title: 'AREA',
      home: LoginPage(),
    );
  }
}
