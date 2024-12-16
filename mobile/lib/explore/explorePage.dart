import 'package:flutter/material.dart';

import "../custom_widget.dart";
import "../footer_page.dart";
import "exploreCardListFiller.dart";
import "topBarExplore.dart";
import "exploreWidget.dart";
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
    if (SelectedCardList.cardList.isEmpty) {
      updateCardList();
    }
  }

  Future<void> updateCardList() async {
    setState(() {
      if (SelectedCardList.index == 0) {
        FillServicesAndWorkflow.fillServicesAndWorkflowList();
        SelectedCardList.cardList = FillServicesAndWorkflow.servicesAndWorkflow;
      } else if (SelectedCardList.index == 1) {
        FillWorkflowsCard.fillWorkflowsCardList();
        SelectedCardList.cardList = FillWorkflowsCard.workflowsCardList;
      } else {
        FillServicesCard.fillServicesCardList();
        SelectedCardList.cardList = FillServicesCard.serviceCardList;
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SingleChildScrollView(
        child: Column(
          children: [
            Padding(padding: EdgeInsets.only(top: 30)),
            TopBarExplore(
              onPressed: updateCardList,
            ),
            MyDivider(),
            MySearchBar(),
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
