import 'package:flutter/material.dart';
import 'package:auto_size_text/auto_size_text.dart';

import "../custom_widget.dart";
import "../footer_page.dart";
import "exploreCardListFiller.dart";
import "explorePageData.dart";

class ExplorePage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      theme: new ThemeData(
          scaffoldBackgroundColor: const Color.fromARGB(255, 255, 255, 255)),
      home: ExplorePageStateful(),
    );
  }
}

class ExplorePageStateful extends StatefulWidget {
  @override
  ExplorePageMain createState() => ExplorePageMain();
}

class ExplorePageMain extends State<ExplorePageStateful> {
  @override
  void initState() {
    super.initState();
    if (SelectedCardList.cardList.isNotEmpty) {
      SelectedCardList.cardList.clear();
    }
    updateCardList();
  }

  Future<void> updateCardList() async {
    setState(() {
      FillServicesCard.fillServicesCardList();
      SelectedCardList.cardList = FillServicesCard.serviceCardList;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SingleChildScrollView(
        child: Column(
          children: [
            Padding(padding: EdgeInsets.only(top: 30)),
            AutoSizeText(
              "Explore",
              style: TextStyle(
                fontSize: 30,
                fontWeight: FontWeight.w900,
                color: Colors.black,
                decorationThickness: 2,
              ),
              maxLines: 1,
            ),
            MyDivider(),
            Column(
              children: SelectedCardList.cardList,
            ),
          ],
        ),
      ),
      bottomNavigationBar: FooterNavigationBar(currentIndex: 1),
    );
  }
}
