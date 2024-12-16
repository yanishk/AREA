import 'package:flutter/material.dart';
import 'dart:developer' as developer;

import "../custom_widget.dart";
import "../footer_page.dart";
import "exploreCardListFiller.dart";
import "topBarExplore.dart";
import "exploreWidget.dart";

class ExplorePage extends StatelessWidget {
  const ExplorePage({
    super.key,
    this.explorePageIndex = 0,
  });
  final int explorePageIndex;
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      theme: new ThemeData(
          scaffoldBackgroundColor: const Color.fromARGB(255, 255, 255, 255)),
      home: ExplorePageStateful(explorePageIndex: explorePageIndex),
    );
  }
}

class ExplorePageStateful extends StatefulWidget {
  const ExplorePageStateful({
    super.key,
    this.explorePageIndex = 0,
  });
  final int explorePageIndex;
  @override
  ExplorePageMain createState() => ExplorePageMain(explorePageIndex: explorePageIndex);
}

class ExplorePageMain extends State<ExplorePageStateful> {
  ExplorePageMain({
    required this.explorePageIndex,
  });
  final int explorePageIndex;

  List<Widget> cardList = [];

  @override
  void initState() {
    super.initState();
    developer.log("explorePageIndex : $explorePageIndex");
    if (explorePageIndex == 0) {
      FillServicesAndWorkflow.fillServicesAndWorkflowList();
      cardList = FillServicesAndWorkflow.servicesAndWorkflow;
    } else if (explorePageIndex == 1) {
      FillWorkflowsCard.fillWorkflowsCardList();
      cardList = FillWorkflowsCard.workflowsCardList;
    } else {
      FillServicesCard.fillServicesCardList();
      cardList = FillServicesCard.serviceCardList;
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SingleChildScrollView(
        child: Column(
          children: [
            Padding(padding: EdgeInsets.only(top: 30)),
            TopBarExplore(
              currentIndex: explorePageIndex,
            ),
            MyDivider(),
            MySearchBar(),
            Column(
              children: cardList,
            ),
          ],
        ),
      ),
      bottomNavigationBar: FooterNavigationBar(currentIndex: 1),
    );
  }
}
