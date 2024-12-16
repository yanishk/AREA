import 'package:flutter/material.dart';

import "../custom_widget.dart";
import "../footer_page.dart";
import "topBarExplore.dart";
import "exploreWidget.dart";

class ExplorePageServices extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      theme: new ThemeData(
          scaffoldBackgroundColor: const Color.fromARGB(255, 255, 255, 255)),
      home: ExplorePageServicesMain(),
    );
  }
}

class ExplorePageServicesMain extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SingleChildScrollView(
        child: Column(
          children: [
            Padding(padding: EdgeInsets.only(top: 30)),
            TopBarExplore(
              currentIndex: 2,
            ),
            MyDivider(),
            MySearchBar(),
            ServicesCard(),
            Padding(padding: EdgeInsets.only(top: 10)),
            ServicesCard(),
            Padding(padding: EdgeInsets.only(top: 10)),
            ServicesCard(),
          ],
        ),
      ),
      bottomNavigationBar: FooterNavigationBar(currentIndex: 1),
    );
  }
}