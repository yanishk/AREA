import 'package:flutter/material.dart';

import "../custom_widget.dart";
import "../footer_page.dart";
import "../globalData.dart";
import "topBarExplore.dart";
import "exploreWidget.dart";

class ExplorePageAll extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      theme: new ThemeData(
          scaffoldBackgroundColor: const Color.fromARGB(255, 255, 255, 255)),
      home: ExplorePageAllStateful(),
    );
  }
}

class ExplorePageAllStateful extends StatefulWidget {
  @override
  ExplorePageAllMain createState() => ExplorePageAllMain();
}

class ExplorePageAllMain extends State<ExplorePageAllStateful> {
  int serviceIndex = 0;
  int allIndex = 0;
  List<String> serviceNameList = Globaldata.actionServicesName + Globaldata.reactionServicesName;
  List<String> serviceColorList = Globaldata.actionServicesColor + Globaldata.reactionServicesColor;
  List<String> serviceIconsList = Globaldata.actionServicesIcons + Globaldata.reactionServicesIcons;

  Widget serviceTiles(int index) {
    if (index % 2 == 0) {
      if (serviceIndex < serviceNameList.length) {
        return ServicesCard(
          name: serviceNameList[serviceIndex],
          color: HexColor(serviceColorList[serviceIndex]),
          imageUrl: serviceIconsList[serviceIndex],
        );
      }
    }
    return WorkflowCard();
  }

  @override
  Widget build(BuildContext context) {
    List<Widget> servicesAndWorkflow = [];
    for (allIndex = 0; allIndex < 5; allIndex++) {
      servicesAndWorkflow.add(serviceTiles(allIndex));
      if (allIndex % 2 == 0) {
        if (serviceIndex < serviceNameList.length) {
          serviceIndex++;
        }
      }
    }
    return Scaffold(
      body: SingleChildScrollView(
        child: Column(
          children: [
            Padding(padding: EdgeInsets.only(top: 30)),
            TopBarExplore(
              currentIndex: 0,
            ),
            MyDivider(),
            MySearchBar(),
            Column(
              children: servicesAndWorkflow,
            ),
          ],
        ),
      ),
      bottomNavigationBar: FooterNavigationBar(currentIndex: 1),
    );
  }
}
