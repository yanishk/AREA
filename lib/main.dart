import 'package:flutter/material.dart';

import "explore.dart";

void main() {
  runApp(MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      debugShowCheckedModeBanner: false,
      title: 'AREA',
      home: ExplorePage(),
    );
  }
}
